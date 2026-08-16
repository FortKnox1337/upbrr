// Copyright (c) 2025-2026, Audionut and the autobrr contributors.
// SPDX-License-Identifier: GPL-2.0-or-later

package huno

import (
	"github.com/autobrr/upbrr/internal/trackers"
	"github.com/autobrr/upbrr/internal/trackers/impl/unit3d"
)

// Profile returns HUNO's Unit3D manifest, including its query-token API and
// multipart text-file upload protocol.
func Profile() unit3d.Profile {
	return unit3d.Profile{
		Name:             "HUNO",
		BaseURL:          "https://hawke.uno",
		Rules:            Rules(),
		AudioPolicy:      AudioPolicy(),
		ValidationPolicy: ValidationPolicy(),
		BannedGroups:     BannedGroups(),
		UploadArtifact: &trackers.UploadArtifactPolicy{
			Source:          "HUNO",
			RequireAnnounce: true,
		},
		ImageHost: &trackers.ImageHostPolicy{
			AllowedHosts: []string{"imgbox", "imgbb", "pixhost", "bam", "onlyimage", "ptscreens", "passtheimage", "hawke.pics"},
		},
		Site: unit3d.SiteProfile{
			BuildDescription: buildDescription,
			OmitNFO:          true,
			UploadAPIKeyTransport: &trackers.APIKeyTransportPolicy{
				QueryParameter: "api_token",
				DisableBearer:  true,
			},
			MultipartFiles: unit3d.MultipartFileProfile{
				DescriptionFilename: "description.txt",
				MediaInfoFilename:   "mediainfo.txt",
				BDInfoFilename:      "bdinfo.txt",
				TorrentFromName:     true,
			},
			PayloadFields: []string{
				"category_id",
				"type_id",
				"tmdb",
				"anonymous",
				"imdb",
				"internal",
				"edition",
				"release_tag",
				"region",
				"distributor",
				"season_number",
				"episode_number",
				"tvdb",
				"mal",
				"season_pack",
				"description",
				"mediainfo",
				"bdinfo",
			},
			ResolveTypeID:          typeID,
			ResolveResolutionID:    resolutionID,
			ApplyAdditionalPayload: applyAdditionalPayload,
		},
	}
}
