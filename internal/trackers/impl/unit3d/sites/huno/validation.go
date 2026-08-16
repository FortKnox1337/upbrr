// Copyright (c) 2025-2026, Audionut and the autobrr contributors.
// SPDX-License-Identifier: GPL-2.0-or-later

package huno

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/autobrr/upbrr/internal/trackers"
	"github.com/autobrr/upbrr/internal/trackers/impl/unit3d"
	"github.com/autobrr/upbrr/pkg/api"
)

const hunoMinimumFallbackBitrate = int64(3_000_000)

var hunoWebRipPattern = regexp.MustCompile(`(?i)(^|[^[:alnum:]])web-?rip([^[:alnum:]]|$)`)

func checkRequirements(ctx context.Context, subject api.TrackerValidationSubject, _ api.Logger) ([]api.RuleFailure, error) {
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("context canceled: %w", err)
	}
	ruleSubject := unit3d.ValidationRuleSubject(subject)
	if hunoWebRip(subject) {
		return []api.RuleFailure{trackers.NewRuleFailure(
			"huno_webrip",
			"WEBRips are not allowed at HUNO.",
			api.RuleDispositionStrict,
		)}, nil
	}
	if !hunoQualityType(unit3d.RuleType(ruleSubject)) {
		return nil, nil
	}
	crf := subject.Assessments.VideoCRF
	if crf.Status == api.VideoCRFStatusPresent {
		if crf.Value <= 22 {
			return nil, nil
		}
		return []api.RuleFailure{trackers.NewRuleFailure(
			"huno_video_crf",
			"HUNO does not allow CRF values above 22.",
			api.RuleDispositionStrict,
		)}, nil
	}
	if unit3d.Animation(ruleSubject) {
		return nil, nil
	}
	bitrate := subject.Assessments.VideoBitrate
	if bitrate.Status != api.VideoBitrateStatusPresent {
		return []api.RuleFailure{trackers.NewRuleFailure(
			"huno_encoding_quality",
			"HUNO requires a valid CRF value or prepared MediaInfo video bitrate.",
			api.RuleDispositionStrict,
		)}, nil
	}
	if bitrate.BitsPerSecond < hunoMinimumFallbackBitrate {
		return []api.RuleFailure{trackers.NewRuleFailure(
			"huno_video_bitrate",
			"Video bitrate is below HUNO's 3 Mbps fallback minimum.",
			api.RuleDispositionStrict,
		)}, nil
	}
	return nil, nil
}

func hunoWebRip(subject api.TrackerValidationSubject) bool {
	return strings.EqualFold(strings.TrimSpace(subject.Type), "WEBRIP") ||
		hunoWebRipPattern.MatchString(subject.Source) ||
		hunoWebRipPattern.MatchString(subject.Release.Source) ||
		hunoWebRipPattern.MatchString(subject.ReleaseName) ||
		hunoWebRipPattern.MatchString(subject.ReleaseNameNoTag)
}

func hunoQualityType(value string) bool {
	switch strings.ToUpper(strings.TrimSpace(value)) {
	case "ENCODE", "DVDRIP", "HDTV":
		return true
	default:
		return false
	}
}
