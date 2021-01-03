// Copyright 2021 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package gitlab

import (
	"net/http"

	"code.gitea.io/gitea/modules/log"
)

func (s *Server) unregisterHadler(w http.ResponseWriter, r *http.Request) {
	log.Warn("Unegister")
	w.WriteHeader(http.StatusOK)
}
