// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package build

import (
	"net/http"

	"github.com/go-chi/chi"
)

// RegisterRoutes registers all routes to handle runner requests.
func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()

	s.router = r

	r.Use(s.authorization)
	r.Put("/runner", s.connectHadler)
	// r.Delete("/runner", s.disconnectHandler)

	return s.router
}
