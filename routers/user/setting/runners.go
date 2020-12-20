// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package setting

import (
	"code.gitea.io/gitea/models"
	"code.gitea.io/gitea/modules/build"
	"code.gitea.io/gitea/modules/context"
	"code.gitea.io/gitea/modules/setting"
)

// GiteaRunnerNewPost creates new Gitea build runner
func GiteaRunnerNewPost(ctx *context.Context) {
	ctx.Data["Title"] = ctx.Tr("settings")
	ctx.Data["PageIsSettingsApplications"] = true
	ctx.Data["BaseLink"] = setting.AppSubURL + "/user/settings"
	ctx.Data["RunnerDescription"] = ctx.Tr("settings.runners_desc")

	if ctx.HasError() {
		loadApplicationsData(ctx)

		ctx.HTML(200, tplSettingsApplications)
		return
	}

	if err := models.NewBuildRunner(ctx.User.ID, models.GiteaRunner); err != nil {
		ctx.ServerError("NewBuildRunner", err)
		return
	}

	ctx.Flash.Success(ctx.Tr("settings.add_runner_success"))
	ctx.Redirect(setting.AppSubURL + "/user/settings/applications")
}

// DeleteRunner response for delete user build runner
func DeleteRunner(ctx *context.Context) {
	if br, err := models.DeleteBuildRunnerByOwnerID(ctx.User.ID, ctx.QueryInt64("id")); err != nil {
		ctx.Flash.Error("DeleteRunnerByOwnerID: " + err.Error())
	} else {
		if (br.Type == models.GiteaRunner) {
			build.InvalidateRunnerBySecret(br.Secret)
		}
		ctx.Flash.Success(ctx.Tr("settings.runner_deletion_success"))
	}

	ctx.JSON(200, map[string]interface{}{
		"redirect": setting.AppSubURL + "/user/settings/applications",
	})
}
