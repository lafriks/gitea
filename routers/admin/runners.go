// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package admin

import (
	"code.gitea.io/gitea/models"
	"code.gitea.io/gitea/modules/base"
	"code.gitea.io/gitea/modules/build"
	"code.gitea.io/gitea/modules/context"
	"code.gitea.io/gitea/modules/setting"
)

const (
	tplAdminRunners base.TplName = "admin/runners"
)

// Runners render runner list page
func Runners(ctx *context.Context) {
	ctx.Data["Title"] = ctx.Tr("admin.runners")
	ctx.Data["PageIsAdminRunners"] = true
	ctx.Data["BaseLink"] = setting.AppSubURL + "/admin"
	ctx.Data["RunnerDescription"] = ctx.Tr("admin.runners.desc")

	loadRunnersData(ctx)

	ctx.HTML(200, tplAdminRunners)
}

// GiteaRunnerNewPost creates new Gitea build runner
func GiteaRunnerNewPost(ctx *context.Context) {
	ctx.Data["Title"] = ctx.Tr("admin.runners")
	ctx.Data["PageIsAdminRunners"] = true
	ctx.Data["BaseLink"] = setting.AppSubURL + "/admin"
	ctx.Data["RunnerDescription"] = ctx.Tr("admin.runners.desc")

	if ctx.HasError() {
		loadRunnersData(ctx)

		ctx.HTML(200, tplAdminRunners)
		return
	}

	if err := models.NewSystemBuildRunner(models.GiteaRunner); err != nil {
		ctx.ServerError("NewSystemBuildRunner", err)
		return
	}

	ctx.Flash.Success(ctx.Tr("admin.runners.add_runner_success"))
	ctx.Redirect(setting.AppSubURL + "/admin/runners")
}

func loadRunnersData(ctx *context.Context) {
	var err error
	ctx.Data["Runners"], err = models.GetSystemBuildRunners(models.ListOptions{})
	if err != nil {
		ctx.ServerError("GetSystemRunners", err)
		return
	}
}

// DeleteRunner response for delete runner
func DeleteRunner(ctx *context.Context) {
	if br, err := models.DeleteSystemBuildRunner(ctx.QueryInt64("id")); err != nil {
		ctx.Flash.Error("DeleteSystemRunner: " + err.Error())
	} else {
		if (br.Type == models.GiteaRunner) {
			build.InvalidateRunnerBySecret(br.Secret)
		}
		ctx.Flash.Success(ctx.Tr("admin.runners.runner_deletion_success"))
	}

	ctx.JSON(200, map[string]interface{}{
		"redirect": setting.AppSubURL + "/admin/runners",
	})
}
