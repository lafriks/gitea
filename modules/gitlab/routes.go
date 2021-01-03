// Copyright 2021 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package gitlab

import (
	"net/http"

	"github.com/go-chi/chi"
)

// RegisterRoutes registers all routes to handle runner requests.
func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()

	s.router = r

	r.Use(s.authorization)
	r.Post("/runners", s.registerHadler)
	r.Delete("/runners", s.unregisterHadler)

	return s.router
}

/*
docker run -d --name gitlab-runner --restart always \
     -v /srv/gitlab-runner/config:/etc/gitlab-runner \
     -v /var/run/docker.sock:/var/run/docker.sock \
		 gitlab/gitlab-runner:latest
*/
