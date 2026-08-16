// Copyright (c) 2025-2026, Audionut and the autobrr contributors.
// SPDX-License-Identifier: GPL-2.0-or-later

package metadata

import (
	"testing"

	"github.com/autobrr/upbrr/pkg/api"
)

func TestVideoCRFAssessment(t *testing.T) {
	tests := []struct {
		name string
		json string
		want api.VideoCRFAssessment
	}{
		{
			name: "present integer",
			json: `{"media":{"track":[{"@type":"Video","Encoded_Library_Settings":"crf=22 / aq-mode=3"}]}}`,
			want: api.VideoCRFAssessment{Status: api.VideoCRFStatusPresent, Value: 22},
		},
		{
			name: "present decimal",
			json: `{"media":{"track":[{"@type":"Video","Encoded_Library_Settings":"cabac=1 / crf = 18.5 / ref=4"}]}}`,
			want: api.VideoCRFAssessment{Status: api.VideoCRFStatusPresent, Value: 18.5},
		},
		{
			name: "two pass unavailable",
			json: `{"media":{"track":[{"@type":"Video","Encoded_Library_Settings":"bitrate=5000 / pass=2"}]}}`,
			want: api.VideoCRFAssessment{Status: api.VideoCRFStatusUnavailable},
		},
		{
			name: "malformed",
			json: `{"media":{"track":[{"@type":"Video","Encoded_Library_Settings":"crf=invalid / aq-mode=3"}]}}`,
			want: api.VideoCRFAssessment{Status: api.VideoCRFStatusInvalid},
		},
		{
			name: "missing settings",
			json: `{"media":{"track":[{"@type":"Video","Format":"HEVC"}]}}`,
			want: api.VideoCRFAssessment{Status: api.VideoCRFStatusUnavailable},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := videoCRFAssessment(mustParseMediaInfoDoc(test.json))
			if got != test.want {
				t.Fatalf("video CRF assessment = %#v, want %#v", got, test.want)
			}
		})
	}
}
