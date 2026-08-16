// Copyright (c) 2025-2026, Audionut and the autobrr contributors.
// SPDX-License-Identifier: GPL-2.0-or-later

package huno

import "github.com/autobrr/upbrr/internal/trackers"

// Rules returns HUNO's waivable language rule and HEVC requirements. Valid
// MediaInfo settings are composed automatically for every Unit3D tracker.
func Rules() *trackers.RuleSet {
	return &trackers.RuleSet{
		RequireAudioLanguages: true,
		RequireHEVCForTypes:   []string{"ENCODE", "DVDRIP", "HDTV"},
	}
}

// AudioPolicy permits additional non-original audio tracks. HUNO explicitly
// allows bloated audio, so the shared language assessment must not warn for it.
func AudioPolicy() *trackers.AudioPolicy {
	return &trackers.AudioPolicy{AllowBloat: true}
}

// ValidationPolicy returns HUNO's WEBRip and encoding-quality requirements.
func ValidationPolicy() trackers.ValidationPolicyBinding {
	return trackers.ValidationPolicyBinding{
		ID:    "unit3d-huno-policy-v1",
		Check: checkRequirements,
	}
}
