// Copyright 2021 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package gitlab

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"code.gitea.io/gitea/modules/log"

	"code.gitea.io/gitea/modules/gitlab/common"
)

func (s *Server) registerHadler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error("Body.ReadAll: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	req := &common.RegisterRunnerRequest{}
	if err = json.Unmarshal(body, req); err != nil {
		log.Error("RegisterRunnerRequest.Unmarshal: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Warn("Register: %s", req.Description)

	res := &common.RegisterRunnerResponse{
		Token: req.Token,
	}

	j, err := json.Marshal(res)
	if err != nil {
		log.Error("RegisterRunnerResponse.Marshal: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(j)
}
