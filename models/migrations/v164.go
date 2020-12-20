// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package migrations

import (
	"code.gitea.io/gitea/modules/timeutil"

	"xorm.io/xorm"
)

func addGiteaBuilds(x *xorm.Engine) error {
	type BuildRunner struct {
		ID          int64              `xorm:"pk autoincr"`
		Type        string             `xorm:"VARCHAR(16) NOT NULL"`
		OwnerID     int64              `xorm:"INDEX"`
		Secret      string             `xorm:"VARCHAR(64) INDEX"`
		CreatedUnix timeutil.TimeStamp `xorm:"created"`
	}

	return x.Sync2(new(BuildRunner))
}
