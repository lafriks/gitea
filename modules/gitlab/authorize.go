// Copyright 2021 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package gitlab

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"code.gitea.io/gitea/models"
	"code.gitea.io/gitea/modules/log"
)

type authorizationRequest struct {
	Token string `json:"Token"`
}

func (s *Server) authorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Job-Token")
		if len(token) == 0 {
			// Read the content
			var bodyBytes []byte
			if r.Body != nil {
				bodyBytes, _ = ioutil.ReadAll(r.Body)
			}
			// Restore the io.ReadCloser to its original state
			r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))// Use the content

			req := &authorizationRequest{}
			if err := json.Unmarshal(bodyBytes, req); err != nil {
				log.Error("Read token from body: %v", err)
			} else {
				token = req.Token
			}
		}

		if len(token) == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Find runner by provided secret in header
		br, err := s.GetRunnerBySecret(token)
		if err != nil {
			if (models.IsErrBuildRunnerNotExist(err)) {
				w.WriteHeader(http.StatusForbidden)
			} else {
				log.Error("GetRunnerBySecret: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}
		if br == nil || br.Type != models.GitLabRunner {
			w.WriteHeader(http.StatusForbidden)
		} else {
			r = r.WithContext(context.WithValue(r.Context(), RunnerCtxKey, br))
			next.ServeHTTP(w, r)
		}
	})
}
