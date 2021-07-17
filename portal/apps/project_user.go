package apps

import (
	"cloudiac/portal/consts/e"
	"cloudiac/portal/libs/ctx"
	"cloudiac/portal/models"
	"cloudiac/portal/models/forms"
	"cloudiac/portal/services"
	"fmt"
	"net/http"
)

func CreateProjectUser(c *ctx.ServiceCtx, form *forms.CreateProjectUserForm) (interface{}, e.Error) {
	pu, err := services.CreateProjectUser(c.DB(), models.UserProject{
		Role:      form.Role,
		UserId:    form.UserId,
		ProjectId: c.ProjectId,
	})
	if err != nil && err.Code() == e.ProjectUserAlreadyExists {
		return nil, e.New(err.Code(), err, http.StatusBadRequest)
	} else if err != nil {
		c.Logger().Errorf("error create project user, err %s", err)
		return nil, e.AutoNew(err, e.DBError)
	}
	return pu, nil
}

// SearchProjectUser 查询组织下某个项目的用户(不包含已经关联项目的用户)
func SearchProjectUser(c *ctx.ServiceCtx) (interface{}, e.Error) {
	query := services.QueryUser(c.DB()).
		Where("status = 'enable'").
		Order("created_at DESC")

	if c.OrgId != "" {
		userIds, _ := services.GetUserIdsByOrg(c.DB(), c.OrgId)
		query = query.Where(fmt.Sprintf("%s.id in (?)", models.User{}.TableName()), userIds)
	}
	if c.ProjectId != "" {
		userIds, _ := services.GetUserIdsByProjectUser(c.DB(), c.ProjectId)
		query = query.Where(fmt.Sprintf("%s.id not in (?)", models.User{}.TableName()), userIds)
	}

	// 导出用户角色

	if c.OrgId != "" {
		query = query.Joins(fmt.Sprintf("left join %s as o on %s.id = o.user_id and o.org_id = ?",
			models.UserOrg{}.TableName(), models.User{}.TableName()), c.OrgId).
			LazySelectAppend(fmt.Sprintf("o.role,%s.*", models.User{}.TableName()))
	}

	users := make([]*models.UserWithRoleResp, 0)
	if err := query.Scan(&users); err != nil {
		c.Logger().Errorf("error get users, err %s", err)
		return nil, e.New(e.DBError, err)
	}

	return users, nil
}

func UpdateProjectUser(c *ctx.ServiceCtx, form *forms.UpdateProjectUserForm) (interface{}, e.Error) {
	if !c.IsSuperAdmin && c.OrgId == "" {
		return nil, e.New(e.PermissionDeny, fmt.Errorf("super admin required"), http.StatusBadRequest)
	}

	attrs := models.Attrs{}
	if form.HasKey("role") {
		attrs["role"] = form.Role
	}

	err := services.UpdateProjectUser(c.DB().
		Where("user_id = ?", form.Id).
		Where("project_id = ?", c.ProjectId),
		form.Id, attrs)
	if err != nil && err.Code() == e.ProjectUserAliasDuplicate {
		return nil, e.New(err.Code(), err, http.StatusBadRequest)
	} else if err != nil {
		c.Logger().Errorf("error update project user, err %s", err)
		return nil, err
	}
	return nil, nil
}

func DeleteProjectUser(c *ctx.ServiceCtx, form *forms.DeleteProjectOrgUserForm) (interface{}, e.Error) {
	return nil, services.DeleteProjectUser(c.DB(), form.Id)
}

func SearchProjectAuthorizationUser(c *ctx.ServiceCtx, form *forms.SearchProjectAuthorizationUserForm) (interface{}, e.Error) {
	query := services.QueryUser(c.DB()).
		Where("status = 'enable'").
		Order("created_at DESC")

	if c.ProjectId != "" {
		userIds, _ := services.GetUserIdsByProjectUser(c.DB(), c.ProjectId)
		query = query.Where(fmt.Sprintf("%s.id  in (?)", models.User{}.TableName()), userIds)
	}

	// 导出用户角色
	if c.OrgId != "" {
		query = query.Joins(fmt.Sprintf("left join %s as o on %s.id = o.user_id and o.org_id = ?",
			models.UserOrg{}.TableName(), models.User{}.TableName()), c.OrgId).
			LazySelectAppend(fmt.Sprintf("o.role,%s.*", models.User{}.TableName()))
	}

	rs, err := getPage(query, form, &models.UserWithRoleResp{})
	if err != nil {
		c.Logger().Errorf("error get page, err %s", err)
	}

	return rs, nil
}
