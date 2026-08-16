// Copyright (c) 2025-2026, Audionut and the autobrr contributors.
// SPDX-License-Identifier: GPL-2.0-or-later

package huno

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/autobrr/upbrr/internal/config"
	descriptionunit3d "github.com/autobrr/upbrr/internal/description/unit3d"
	"github.com/autobrr/upbrr/pkg/api"
)

const hunoMaxScreenshotRowWidth = 1100
const hunoDefaultScreenshotWidth = 350
const hunoDefaultScreensPerRow = 2
const hunoSignatureSize = 8

func buildDescription(
	ctx context.Context,
	meta api.UploadSubject,
	appConfig config.Config,
	trackerConfig config.TrackerConfig,
	logger api.Logger,
	keptDescription string,
	menuImages []api.ScreenshotImage,
	screenshots []api.ScreenshotImage,
) (string, error) {
	cfg := appConfig
	thumbnailWidth := cfg.Description.ThumbnailSize
	if thumbnailWidth <= 0 {
		thumbnailWidth = hunoDefaultScreenshotWidth
	}
	screensPerRow := hunoDefaultScreensPerRow
	if configured, parseErr := strconv.Atoi(strings.TrimSpace(cfg.Description.ScreensPerRow)); parseErr == nil && configured > 0 {
		screensPerRow = configured
	}
	for screensPerRow > 1 && screensPerRow*thumbnailWidth > hunoMaxScreenshotRowWidth {
		screensPerRow--
	}
	cfg.Description.ScreensPerRow = strconv.Itoa(screensPerRow)

	value, err := descriptionunit3d.BuildDescription(
		ctx,
		api.NewDescriptionSubject(meta),
		cfg,
		trackerConfig,
		logger,
		keptDescription,
		menuImages,
		screenshots,
	)
	if err != nil {
		return "", fmt.Errorf("trackers: %w", err)
	}
	link, text := descriptionunit3d.UppbrrSignatureLink()
	defaultSignature := fmt.Sprintf("[right][url=%s][size=4]%s[/size][/url][/right]", link, text)
	hunoSignature := fmt.Sprintf("[right][url=%s][size=%d]%s[/size][/url][/right]", link, hunoSignatureSize, text)
	value = strings.Replace(value, defaultSignature, hunoSignature, 1)
	if strings.Contains(value, hunoSignature) {
		return value, nil
	}
	if strings.TrimSpace(value) == "" {
		return hunoSignature, nil
	}
	return strings.TrimSpace(value) + "\n\n" + hunoSignature, nil
}
