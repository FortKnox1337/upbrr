// Copyright (c) 2025-2026, Audionut and the autobrr contributors.
// SPDX-License-Identifier: GPL-2.0-or-later

package huno

import (
	"strings"

	"github.com/autobrr/upbrr/internal/trackers"
)

func applyAdditionalPayload(req trackers.PreparationInput, data map[string]string) {
	meta := req.Meta
	if !req.Runtime.Internal {
		delete(data, "internal")
	}
	if edition := strings.TrimSpace(meta.Edition); edition != "" {
		data["edition"] = edition
	}
	if releaseTag := strings.TrimSpace(meta.Repack); releaseTag != "" {
		data["release_tag"] = releaseTag
	}
	if strings.TrimSpace(meta.DiscType) != "" {
		if region := strings.TrimSpace(meta.Region); region != "" {
			data["region"] = region
		}
		if distributor := strings.TrimSpace(meta.Distributor); distributor != "" {
			data["distributor"] = distributor
		}
	}
	if strings.EqualFold(strings.TrimSpace(string(meta.Identity.Category)), "tv") {
		data["season_pack"] = boolString(meta.TVPack)
	} else {
		delete(data, "mal")
	}
}

func boolString(value bool) string {
	if value {
		return "1"
	}
	return "0"
}
