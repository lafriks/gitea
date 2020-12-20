// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package build

import (
	"context"
	"net/http"

	"code.gitea.io/gitea/models"
	"code.gitea.io/gitea/modules/log"
)

func (s *Server) authorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Find runner by provided secret in header
		br, err := s.GetRunnerBySecret(r.Header.Get("X-Authorization-Token"))
		if err != nil {
			if (models.IsErrBuildRunnerNotExist(err)) {
				w.WriteHeader(http.StatusUnauthorized)
			} else {
				log.Error("GetRunnerBySecret: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}
		if br == nil || br.Type != models.GiteaRunner {
			w.WriteHeader(403)
		} else {
			r = r.WithContext(context.WithValue(r.Context(), RunnerCtxKey, br))
			next.ServeHTTP(w, r)
		}
	})
}
