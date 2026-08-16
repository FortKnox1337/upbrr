// Copyright (c) 2025-2026, Audionut and the autobrr contributors.
// SPDX-License-Identifier: GPL-2.0-or-later

package api

import "testing"

func TestNewTrackerValidationSubjectProjectsVideoCRFAssessment(t *testing.T) {
	assessment := VideoCRFAssessment{Status: VideoCRFStatusPresent, Value: 18.5}
	subject := NewTrackerValidationSubject(UploadSubject{Assessments: ReleaseAssessments{VideoCRF: assessment}}, "EXAMPLE")
	if subject.Assessments.VideoCRF != assessment {
		t.Fatalf("projected video CRF = %#v, want %#v", subject.Assessments.VideoCRF, assessment)
	}
}
