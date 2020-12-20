// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package build

import (
	"context"
	"net/http"
)

// server represents a Build server instance
var server *Server

// RegisterRoutes registers all routes used by runners.
func RegisterRoutes() http.Handler {
	return server.RegisterRoutes()
}

// InvalidateRunnerBySecret invalidates cached secret
func InvalidateRunnerBySecret(secret string) {
	server.InvalidateRunnerBySecret(secret)
}

// Init initializes the Build server instance
func Init(ctx context.Context) error {
	server = NewServer(ctx)

	return nil
}
