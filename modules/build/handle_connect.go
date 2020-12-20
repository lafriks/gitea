// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package build

import (
	"net/http"

	"code.gitea.io/gitea/modules/log"
)

func (s *Server) connectHadler(w http.ResponseWriter, r *http.Request) {
	log.Warn("Connect")
	w.WriteHeader(http.StatusOK)
}
