// Copyright 2020 The Okteto Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package login

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
)

// Handler handles the authentication using a browser
type Handler struct {
	ctx      context.Context
	state    string
	baseURL  string
	port     int
	response chan string
	errChan  chan error
}

func (a *Handler) handle() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		s := r.URL.Query().Get("state")

		if a.state != s {
			a.errChan <- fmt.Errorf("Invalid request state")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if _, err := w.Write(loginHTML); err != nil {
			a.errChan <- fmt.Errorf("Failed to write to the response: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
		}

		a.response <- code
	}

	return http.HandlerFunc(fn)
}

// AuthorizationURL returns the authorization URL used for login
func (h *Handler) AuthorizationURL() string {
	redirectURL := fmt.Sprintf("http://127.0.0.1:%d/authorization-code/callback?state=%s", h.port, h.state)
	params := url.Values{}
	params.Add("state", h.state)
	params.Add("redirect", redirectURL)

	authorizationURL, err := url.Parse(fmt.Sprintf("%s/auth/authorization-code", h.baseURL))
	if err != nil {
		panic(fmt.Sprintf("failed to build authorizationURL: %s", err))
	}
	authorizationURL.RawQuery = params.Encode()
	return authorizationURL.String()
}

func randToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(b), nil
}
