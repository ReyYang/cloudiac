package handlers

import (
	"cloudiac/apps"
	"cloudiac/consts/e"
	"cloudiac/libs/ctx"
	"cloudiac/models/forms"
	"cloudiac/services"
)

type GitLab struct {}


func (GitLab) ListRepos(c *ctx.GinRequestCtx) {
	form := forms.GetGitProjectsForm{}
	if err := c.Bind(&form); err != nil {
		return
	}
	vcs, err := services.QueryVcsByVcsId(form.VcsId, c.ServiceCtx().Tx())
	if err != nil {
		c.JSONResult(nil,e.New(e.DBError, err))
		return
	}
	if vcs.VcsType == "gitlab"{
		c.JSONResult(apps.ListOrganizationRepos(vcs, &form))
	} else if vcs.VcsType == "gitea" {
		c.JSONResult(apps.ListGiteaOrganizationRepos(vcs, &form))
	}

}


func (GitLab) ListBranches(c *ctx.GinRequestCtx) {
	form := forms.GetGitBranchesForm{}
	if err := c.Bind(&form); err != nil {
		return
	}
	vcs, err := services.QueryVcsByVcsId(form.VcsId, c.ServiceCtx().Tx())
	if err != nil {
		c.JSONResult(nil,e.New(e.DBError, err))
		return
	}
	if vcs.VcsType == "gitlab" {
		c.JSONResult(apps.ListRepositoryBranches(vcs, &form))
	} else if vcs.VcsType == "gitea" {
		c.JSONResult(apps.ListGiteaRepoBranches(vcs, &form))
	}

}

func (GitLab) GetReadmeContent(c *ctx.GinRequestCtx) {
	form := forms.GetReadmeForm{}
	if err := c.Bind(&form); err != nil {
		return
	}
	vcs, err := services.QueryVcsByVcsId(form.VcsId, c.ServiceCtx().Tx())
	if err != nil {
		c.JSONResult(nil,e.New(e.DBError, err))
		return
	}
	if vcs.VcsType == "gitlab" {
		c.JSONResult(apps.GetReadmeContent(vcs, &form))
	} else if vcs.VcsType == "gitea" {
		c.JSONResult(apps.GetGiteaReadme(vcs, &form))
	}

}

func TemplateTfvarsSearch(c *ctx.GinRequestCtx){
	form := forms.TemplateTfvarsSearchForm{}
	if err := c.Bind(&form); err != nil {
		return
	}
	vcs, err := services.QueryVcsByVcsId(form.VcsId, c.ServiceCtx().Tx())
	if err != nil {
		c.JSONResult(nil,e.New(e.DBError, err))
		return
	}
	c.JSONResult(apps.TemplateTfvarsSearch(vcs, &form))
}
