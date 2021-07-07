package apps

import (
	"cloudiac/portal/consts/e"
	"cloudiac/portal/libs/ctx"
	"cloudiac/portal/models"
	"cloudiac/portal/models/forms"
	"cloudiac/portal/services"
	"cloudiac/utils"
	"fmt"
)

var (
	emailSubjectResetPassword = "重置密码"
	emailBodyResetPassword    = "尊敬的 {{.Name}}：\n\n您的密码已经被重置，这是您的新密码：\n\n密码：\t{{.InitPass}}\n\n请使用新密码登陆系统。\n\n为了保障您的安全，请立即登陆您的账号并修改密码。"
)

func SearchToken(c *ctx.ServiceCtx, form *forms.SearchTokenForm) (interface{}, e.Error) {
	query := services.QueryToken(c.DB())
	query = query.Where("user_id = ?", c.UserId)
	if form.Status != "" {
		query = query.Where("status = ?", form.Status)
	}
	if form.Q != "" {
		qs := "%" + form.Q + "%"
		query = query.Where("description LIKE ?", qs)
	}

	query = query.Order("created_at DESC")
	rs, _ := getPage(query, form, models.Token{})
	return rs, nil
	//p := page.New(form.CurrentPage(), form.PageSize(), query)
	//tokens := make([]*models.Token, 0)
	//if err := p.Scan(&tokens); err != nil {
	//	return nil, e.New(e.DBError, err)
	//}
	//
	//return page.PageResp{
	//	Total:    p.MustTotal(),
	//	PageSize: p.Size,
	//	List:     tokens,
	//}, nil
}

func CreateToken(c *ctx.ServiceCtx, form *forms.CreateTokenForm) (interface{}, e.Error) {
	c.AddLogField("action", fmt.Sprintf("create token for user %s", c.UserId))

	tokenStr := utils.GenGuid("")
	token, err := services.CreateToken(c.DB(), models.Token{
		Description: form.Description,
		UserId:      c.UserId,
		Token:       tokenStr,
	})
	if err != nil {
		return nil, e.AutoNew(err, e.DBError)
	}
	return token, nil
}

func UpdateToken(c *ctx.ServiceCtx, form *forms.UpdateTokenForm) (token *models.Token, err e.Error) {
	c.AddLogField("action", fmt.Sprintf("update token %d", form.Id))
	if form.Id == "" {
		return nil, e.New(e.BadRequest, fmt.Errorf("missing 'id'"))
	}

	attrs := models.Attrs{}
	if form.HasKey("status") {
		attrs["status"] = form.Status
	}

	if form.HasKey("description") {
		attrs["description"] = form.Description
	}

	token, err = services.UpdateToken(c.DB(), form.Id, attrs)
	return
}

func DeleteToken(c *ctx.ServiceCtx, form *forms.DeleteTokenForm) (result interface{}, re e.Error) {
	c.AddLogField("action", fmt.Sprintf("delete token %d", form.Id))
	if err := services.DeleteToken(c.DB(), form.Id); err != nil {
		return nil, err
	}

	return
}
