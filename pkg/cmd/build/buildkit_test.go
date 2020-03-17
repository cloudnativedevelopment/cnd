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

package build

import "testing"

func Test_translateCacheHandler(t *testing.T) {
	var tests = []struct {
		name     string
		input    string
		userID   string
		expected string
	}{
		{
			name:     "no-matched",
			input:    "RUN go build",
			userID:   "userid",
			expected: "RUN go build",
		},
		{
			name:     "matched-id-first",
			input:    "RUN --mount=id=1,type=cache,target=/root/.cache/go-build go build",
			userID:   "userid",
			expected: "RUN --mount=id=userid-1,type=cache,target=/root/.cache/go-build go build",
		},
		{
			name:     "matched-id-last",
			input:    "RUN --mount=type=cache,target=/root/.cache/go-build,id=1 go build",
			userID:   "userid",
			expected: "RUN --mount=type=cache,target=/root/.cache/go-build,id=userid-1 go build",
		},
		{
			name:     "matched-noid",
			input:    "RUN --mount=type=cache,target=/root/.cache/go-build go build",
			userID:   "userid",
			expected: "RUN --mount=id=userid,type=cache,target=/root/.cache/go-build go build",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := translateCacheHandler(tt.input, tt.userID)
			if tt.expected != result {
				t.Errorf("expected %s got %s in test %s", tt.expected, result, tt.name)
			}
		})
	}

}
