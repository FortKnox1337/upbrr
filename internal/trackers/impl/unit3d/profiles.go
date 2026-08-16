// Copyright (c) 2025-2026, Audionut and the autobrr contributors.
// SPDX-License-Identifier: GPL-2.0-or-later

package unit3d

import (
	"context"

	"github.com/autobrr/upbrr/internal/config"
	"github.com/autobrr/upbrr/internal/trackers"
	"github.com/autobrr/upbrr/pkg/api"
)

// SiteProfile contains optional site-owned Unit3D payload callbacks.
type SiteProfile struct {
	// APIKeyTransport overrides bearer authentication for non-standard Unit3D
	// data lookup and duplicate-search APIs.
	APIKeyTransport *trackers.APIKeyTransportPolicy
	// UploadAPIKeyTransport overrides bearer authentication for a site's upload API.
	UploadAPIKeyTransport *trackers.APIKeyTransportPolicy
	// MultipartFiles moves selected text payload fields into multipart file parts.
	MultipartFiles MultipartFileProfile
	// PayloadFields optionally limits the fields submitted to a site's API.
	PayloadFields []string
	// BuildName optionally overrides the generic Unit3D release-name builder.
	BuildName func(meta api.UploadSubject, cfg config.TrackerConfig) string
	// BuildNameVersion identifies a custom BuildName implementation.
	BuildNameVersion string
	// BuildDescription optionally renders site-specific tracker markup.
	BuildDescription func(ctx context.Context, meta api.UploadSubject, appConfig config.Config, trackerConfig config.TrackerConfig, logger api.Logger, keptDescription string, menuImages []api.ScreenshotImage, screenshots []api.ScreenshotImage) (string, error)
	// OmitNFO prevents the generic Unit3D uploader from attaching an optional
	// NFO part when the site's documented multipart contract does not accept it.
	OmitNFO bool
	// ResolveKeywords optionally maps prepared metadata to Unit3D keywords.
	ResolveKeywords func(meta api.UploadSubject) string
	// ResolveTypeID optionally maps prepared metadata to a site type identifier.
	ResolveTypeID func(meta api.UploadSubject) string
	// ResolveResolutionID optionally maps prepared metadata to a site resolution identifier.
	ResolveResolutionID func(meta api.UploadSubject) string
	// ResolveCategoryID optionally maps prepared metadata to a site category identifier.
	ResolveCategoryID func(meta api.UploadSubject) string
	// ApplyAdditionalPayload appends site-owned fields to a prepared payload.
	ApplyAdditionalPayload func(req trackers.PreparationInput, data map[string]string)
	// FinalizeDescription applies final site-owned description transformations.
	FinalizeDescription func(description string, meta api.UploadSubject) string
}

// MultipartFileProfile declares text payload fields that a site requires as
// multipart file parts. Empty filenames leave the standard form fields intact.
type MultipartFileProfile struct {
	DescriptionFilename string
	MediaInfoFilename   string
	BDInfoFilename      string
	TorrentFromName     bool
}

func firstSiteProfile(profiles []SiteProfile) SiteProfile {
	if len(profiles) == 0 {
		return SiteProfile{}
	}
	return profiles[0]
}
