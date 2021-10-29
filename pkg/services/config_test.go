// Copyright 2021 Paul Greenberg greenpau@outlook.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package services

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestValidate(t *testing.T) {
	testcases := []struct {
		name      string
		units     []*Unit
		want      []*Service
		shouldErr bool
		err       error
	}{
		{
			name: "test reorder config units",
			units: []*Unit{
				{
					Name:    "hostname",
					Command: "hostname",
					Kind:    "command",
				},
				{
					Name:      "test-py-http-server",
					Command:   "python3",
					Kind:      "app",
					Arguments: []string{"-m", "http.server", "4080"},
					After:     []string{"hostname"},
				},
				{
					Name:      "test-py-http-server-4081",
					Kind:      "app",
					Command:   "python3",
					Arguments: []string{"-m", "http.server", "4081"},
					Priority:  100,
				},
			},
			want: []*Service{
				{
					Unit: &Unit{
						Name:    "hostname",
						Command: "hostname",
						Kind:    "command",
					},
					Seq: 0,
				},
				{
					Unit: &Unit{
						Name:      "test-py-http-server",
						Command:   "python3",
						Kind:      "app",
						Arguments: []string{"-m", "http.server", "4080"},
						After:     []string{"hostname"},
					},
					Seq: 1,
				},
				{
					Unit: &Unit{
						Name:      "test-py-http-server-4081",
						Command:   "python3",
						Kind:      "app",
						Arguments: []string{"-m", "http.server", "4081"},
						Priority:  100,
					},
					Seq: 2,
				},
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := NewConfig()
			for _, u := range tc.units {
				if err := cfg.AddUnit(u); err != nil {
					t.Fatal(err)
				}
			}
			got, err := cfg.Services()
			if err != nil {
				if !tc.shouldErr {
					t.Fatalf("expected success, got: %v", err)
				}
				if diff := cmp.Diff(err.Error(), tc.err.Error()); diff != "" {
					t.Fatalf("unexpected error: %v, want: %v", err, tc.err)
				}
				return
			}
			if tc.shouldErr {
				t.Fatalf("unexpected success, want: %v", tc.err)
			}
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("config mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
