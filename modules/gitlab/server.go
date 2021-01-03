// Copyright 2021 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package gitlab

import (
	"context"
	"sync"

	"code.gitea.io/gitea/models"

	"github.com/go-chi/chi"
)

type contextKey struct {
	name string
}

func (k *contextKey) String() string {
	return "runner context value " + k.name
}

var (
	// RunnerCtxKey is the context.Context key to store the request context.
	RunnerCtxKey = &contextKey{"RouteContext"}
)

// Server represent Build server instance
type Server struct {
	ctx    context.Context
	router *chi.Mux
	runnerCache map[string]*models.BuildRunner
	runnerCacheLock *sync.RWMutex
}

// NewServer returns new Build server instance
func NewServer(ctx context.Context) *Server {
	return &Server{
		ctx: ctx,
		runnerCache: make(map[string]*models.BuildRunner),
		runnerCacheLock: &sync.RWMutex{},
	}
}

// GetRunnerBySecret returns runner information either from database or cache
func (s *Server) GetRunnerBySecret(secret string) (*models.BuildRunner, error) {
	s.runnerCacheLock.RLock()

	if br, has := s.runnerCache[secret]; has {
		s.runnerCacheLock.RUnlock()
		return br, nil
	}

	br, err := models.GetBuildRunnerBySecret(secret)
	if err != nil || br == nil {
		s.runnerCacheLock.RUnlock()
		return nil, err
	}

	s.runnerCacheLock.RUnlock()
	s.runnerCacheLock.Lock()
	defer s.runnerCacheLock.Unlock()
	s.runnerCache[secret] = br

	return br, nil
}

// InvalidateRunnerBySecret invalidates cached secret
func (s *Server) InvalidateRunnerBySecret(secret string) {
	s.runnerCacheLock.Lock()
	defer s.runnerCacheLock.Unlock()
	if _, has := s.runnerCache[secret]; has {
		s.runnerCache[secret] = nil
	}
}
