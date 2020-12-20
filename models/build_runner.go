// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package models

import (
	"strings"

	"code.gitea.io/gitea/modules/generate"
	"code.gitea.io/gitea/modules/setting"
	"code.gitea.io/gitea/modules/timeutil"
)

// BuildRunnerType is the type of an build runner
type BuildRunnerType string

// Types of build runner
const (
	GiteaRunner BuildRunnerType = "gitea"
)

// Name returns build runner type display name.
func (brt BuildRunnerType) Name() string {
	switch brt {
	case GiteaRunner:
		return "Gitea"
	default:
		return "Unknown"
	}
}

// BuildRunner represents a instances of build runners
type BuildRunner struct {
	ID          int64              `xorm:"pk autoincr"`
	Type        BuildRunnerType    `xorm:"VARCHAR(16) NOT NULL"`
	OwnerID     int64              `xorm:"INDEX"`
	Secret      string             `xorm:"VARCHAR(64) INDEX"`
	CreatedUnix timeutil.TimeStamp `xorm:"created"`
}

// Settings returns build runner configuration settings
func (br *BuildRunner) Settings() map[string]string {
	r := make(map[string]string, 3)

	if br.Type == GiteaRunner {
		r["GITEA_URL"] = setting.AppURL
		r["GITEA_TOKEN"] = br.Secret
	}

	return r
}

// NewBuildRunner creates new build runner
func NewBuildRunner(ownerID int64, runnerType BuildRunnerType) error {
	secret, err := generate.NewSecretKey()
	if err != nil {
		return err
	}

	br := &BuildRunner{
		Type:    runnerType,
		OwnerID: ownerID,
		Secret:  secret,
	}

	if _, err = x.Insert(br); err != nil {
		return err
	}

	return nil
}

// NewSystemBuildRunner creates new system build runner
func NewSystemBuildRunner(runnerType BuildRunnerType) error {
	secret, err := generate.NewSecretKey()
	if err != nil {
		return err
	}

	br := &BuildRunner{
		Type:    runnerType,
		Secret:  secret,
	}

	if _, err = x.Insert(br); err != nil {
		return err
	}

	return nil
}

// GetBuildRunnersByOwnerID returns paginated runners for an organization or user.
func GetBuildRunnersByOwnerID(ownerID int64, listOptions ListOptions) ([]*BuildRunner, error) {
	if listOptions.Page == 0 {
		br := make([]*BuildRunner, 0, 5)
		return br, x.Find(&br, &BuildRunner{OwnerID: ownerID})
	}

	sess := listOptions.getPaginatedSession()
	br := make([]*BuildRunner, 0, listOptions.PageSize)
	return br, sess.Find(&br, &BuildRunner{OwnerID: ownerID})
}

// GetSystemBuildRunners returns paginated system build runners.
func GetSystemBuildRunners(listOptions ListOptions) ([]*BuildRunner, error) {
	if listOptions.Page == 0 {
		br := make([]*BuildRunner, 0, 5)
		return br, x.Where("`owner_id` = 0").Find(&br)
	}

	sess := listOptions.getPaginatedSession()
	br := make([]*BuildRunner, 0, listOptions.PageSize)
	return br, sess.Where("`owner_id` = 0").Find(&br)
}

// GetBuildRunnerBySecret returns build runner by secret.
func GetBuildRunnerBySecret(secret string) (*BuildRunner, error) {
	br := &BuildRunner{
		Secret: strings.TrimSpace(secret),
	}

	if has, err := x.Get(br); !has {
		return nil, ErrBuildRunnerNotExist{}
	} else if err != nil {
		return nil, err
	}
	return br, nil
}

// deleteBuildRunner uses argument bean as query condition,
// ID must be specified and do not assign unnecessary fields.
func deleteBuildRunner(bean *BuildRunner) (*BuildRunner, error) {
	sess := x.NewSession()
	defer sess.Close()
	if err := sess.Begin(); err != nil {
		return nil, err
	}

	if has, err := sess.Get(bean); !has || err != nil {
		return nil, err
	}

	if count, err := sess.Delete(bean); err != nil {
		return nil, err
	} else if count == 0 {
		return nil, ErrBuildRunnerNotExist{ID: bean.ID}
	}

	return bean, sess.Commit()
}

// DeleteBuildRunnerByOwnerID deletes build runner by given owner ID.
func DeleteBuildRunnerByOwnerID(ownerID, id int64) (*BuildRunner, error) {
	return deleteBuildRunner(&BuildRunner{
		ID:      id,
		OwnerID: ownerID,
	})
}

// DeleteSystemBuildRunner deletes system build runner by given owner ID.
func DeleteSystemBuildRunner(id int64) (*BuildRunner, error) {
	return deleteBuildRunner(&BuildRunner{
		ID:      id,
	})
}
