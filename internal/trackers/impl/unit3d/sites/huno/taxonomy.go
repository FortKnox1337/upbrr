// Copyright (c) 2025-2026, Audionut and the autobrr contributors.
// SPDX-License-Identifier: GPL-2.0-or-later

package huno

import (
	"strings"

	"github.com/autobrr/upbrr/internal/trackers/impl/unit3d"
	"github.com/autobrr/upbrr/pkg/api"
)

func typeID(meta api.UploadSubject) string {
	switch strings.ToUpper(strings.TrimSpace(unit3d.InferType(meta))) {
	case "DISC":
		return "1"
	case "REMUX":
		return "2"
	case "WEBDL":
		return "3"
	case "WEBRIP", "ENCODE", "DVDRIP", "HDTV":
		return "15"
	default:
		return "0"
	}
}

func resolutionID(meta api.UploadSubject) string {
	resolution := strings.ToLower(strings.TrimSpace(unit3d.Resolution(meta)))
	if value, ok := map[string]string{
		"4320p": "1",
		"2160p": "2",
		"1080p": "3",
		"1080i": "4",
		"720p":  "5",
		"576p":  "6",
		"576i":  "7",
		"480p":  "8",
		"480i":  "9",
		"540p":  "11",
		"540i":  "11",
	}[resolution]; ok {
		return value
	}
	return "10"
}
