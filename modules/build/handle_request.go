// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package build

import (
	"context"
	"net/http"
	"time"
)

var defaultTimeout = time.Second * 30

func (s *Server) requestHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	w.WriteHeader(http.StatusOK)
}
