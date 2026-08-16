// Copyright (c) 2025-2026, Audionut and the autobrr contributors.
// SPDX-License-Identifier: GPL-2.0-or-later

package huno

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/autobrr/upbrr/internal/config"
	"github.com/autobrr/upbrr/internal/trackers"
	"github.com/autobrr/upbrr/internal/trackers/impl/unit3d"
	"github.com/autobrr/upbrr/pkg/api"
)

func TestProfile(t *testing.T) {
	profile := Profile()
	if profile.Name != "HUNO" || profile.BaseURL != "https://hawke.uno" {
		t.Fatalf("profile identity = %q/%q", profile.Name, profile.BaseURL)
	}
	if profile.UploadArtifact == nil || profile.UploadArtifact.Source != "HUNO" || !profile.UploadArtifact.RequireAnnounce {
		t.Fatalf("upload artifact policy = %#v", profile.UploadArtifact)
	}
	if profile.Site.APIKeyTransport != nil {
		t.Fatalf("data API key transport = %#v; HUNO data and search APIs use bearer authentication", profile.Site.APIKeyTransport)
	}
	if profile.Site.UploadAPIKeyTransport == nil || profile.Site.UploadAPIKeyTransport.QueryParameter != "api_token" ||
		!profile.Site.UploadAPIKeyTransport.DisableBearer {
		t.Fatalf("upload API key transport = %#v", profile.Site.UploadAPIKeyTransport)
	}
	if profile.Site.MultipartFiles.DescriptionFilename != "description.txt" ||
		profile.Site.MultipartFiles.MediaInfoFilename != "mediainfo.txt" ||
		profile.Site.MultipartFiles.BDInfoFilename != "bdinfo.txt" ||
		!profile.Site.MultipartFiles.TorrentFromName {
		t.Fatalf("multipart text files = %#v", profile.Site.MultipartFiles)
	}
	if profile.Site.BuildDescription == nil {
		t.Fatal("HUNO description builder is not configured")
	}
	if !profile.Site.OmitNFO {
		t.Fatal("HUNO must omit undocumented NFO multipart parts")
	}
	for _, field := range []string{"description", "mediainfo", "bdinfo", "category_id", "type_id"} {
		if !slices.Contains(profile.Site.PayloadFields, field) {
			t.Fatalf("HUNO payload allowlist omitted %q", field)
		}
	}
	for _, field := range []string{"name", "resolution_id", "mod_queue_opt_in"} {
		if slices.Contains(profile.Site.PayloadFields, field) {
			t.Fatalf("HUNO payload allowlist unexpectedly contains %q", field)
		}
	}
	if profile.ImageHost == nil || !slices.Contains(profile.ImageHost.AllowedHosts, "hawke.pics") {
		t.Fatalf("image host policy = %#v", profile.ImageHost)
	}
	if !slices.Contains(profile.BannedGroups, "RARBG") {
		t.Fatal("HUNO banned groups omitted RARBG")
	}

	definition := unit3d.NewWithProfile(profile)
	if definition.TrackerFamily() != trackers.FamilyUnit3D {
		t.Fatalf("tracker family = %q", definition.TrackerFamily())
	}
	if policy := definition.APIKeyTransportPolicy(); policy != nil {
		t.Fatalf("effective data API key transport = %#v; expected default bearer authentication", policy)
	}
	if rules := definition.Rules(); !rules.RequireValidMISetting || !rules.RequireAudioLanguages ||
		!slices.Equal(rules.RequireHEVCForTypes, []string{"ENCODE", "DVDRIP", "HDTV"}) {
		t.Fatalf("effective HUNO rules = %#v", rules)
	}
	if policy := definition.AudioPolicy(); policy == nil || !policy.AllowBloat {
		t.Fatalf("effective HUNO audio policy = %#v", policy)
	}
}

func TestDescriptionUsesConfiguredSizeAndClampsRowsToHUNOMaxWidth(t *testing.T) {
	t.Parallel()

	profile := Profile()
	appConfig := config.Config{}
	appConfig.Description.ThumbnailSize = 500
	appConfig.Description.ScreensPerRow = "3"
	screenshots := make([]api.ScreenshotImage, 0, 6)
	for i := 1; i <= 6; i++ {
		screenshots = append(screenshots, api.ScreenshotImage{
			WebURL: fmt.Sprintf("https://images.example/view/%d", i),
			RawURL: fmt.Sprintf("https://images.example/raw/%d.png", i),
		})
	}
	description, err := profile.Site.BuildDescription(
		context.Background(),
		api.UploadSubject{},
		appConfig,
		config.TrackerConfig{},
		api.NopLogger{},
		"",
		nil,
		screenshots,
	)
	if err != nil {
		t.Fatalf("build HUNO description: %v", err)
	}
	if count := strings.Count(description, "[img=500]"); count != 6 {
		t.Fatalf("configured 500px screenshot count = %d; description=%q", count, description)
	}
	for _, row := range []string{
		"raw/1.png[/img][/url] [url=https://images.example/view/2]",
		"raw/3.png[/img][/url] [url=https://images.example/view/4]",
		"raw/5.png[/img][/url] [url=https://images.example/view/6]",
	} {
		if !strings.Contains(description, row) {
			t.Fatalf("HUNO screenshot row %q missing; description=%q", row, description)
		}
	}
	if !strings.Contains(description, "[size=8]Uploaded by upbrr[/size]") {
		t.Fatalf("HUNO signature size missing; description=%q", description)
	}
}

func TestDescriptionKeepsConfiguredColumnsWhenTheyFit(t *testing.T) {
	t.Parallel()

	profile := Profile()
	appConfig := config.Config{}
	appConfig.Description.ThumbnailSize = 350
	appConfig.Description.ScreensPerRow = "3"
	screenshots := []api.ScreenshotImage{
		{WebURL: "https://images.example/view/1", RawURL: "https://images.example/raw/1.png"},
		{WebURL: "https://images.example/view/2", RawURL: "https://images.example/raw/2.png"},
		{WebURL: "https://images.example/view/3", RawURL: "https://images.example/raw/3.png"},
	}
	description, err := profile.Site.BuildDescription(
		context.Background(),
		api.UploadSubject{},
		appConfig,
		config.TrackerConfig{},
		api.NopLogger{},
		"",
		nil,
		screenshots,
	)
	if err != nil {
		t.Fatalf("build HUNO description: %v", err)
	}
	if count := strings.Count(description, "[img=350]"); count != 3 {
		t.Fatalf("configured 350px screenshot count = %d; description=%q", count, description)
	}
	if !strings.Contains(description, "raw/2.png[/img][/url] [url=https://images.example/view/3]") {
		t.Fatalf("configured three-column row was not preserved; description=%q", description)
	}
}

func TestDescriptionKeepsCustomSignatureAndAppendsHUNOAttribution(t *testing.T) {
	t.Parallel()

	profile := Profile()
	appConfig := config.Config{}
	appConfig.Description.CustomSignature = "[center]Synthetic custom signature[/center]"
	description, err := profile.Site.BuildDescription(
		context.Background(),
		api.UploadSubject{},
		appConfig,
		config.TrackerConfig{},
		api.NopLogger{},
		"",
		nil,
		nil,
	)
	if err != nil {
		t.Fatalf("build HUNO description: %v", err)
	}
	if !strings.Contains(description, appConfig.Description.CustomSignature) {
		t.Fatalf("custom signature missing; description=%q", description)
	}
	if count := strings.Count(description, "[size=8]Uploaded by upbrr[/size]"); count != 1 {
		t.Fatalf("HUNO attribution count = %d; description=%q", count, description)
	}
}

func TestTaxonomy(t *testing.T) {
	profile := Profile()
	typeTests := []struct {
		name string
		meta api.UploadSubject
		want string
	}{
		{
			name: "disc",
			meta: api.UploadSubject{DiscType: "BDMV", Type: "DISC"},
			want: "1",
		},
		{
			name: "remux",
			meta: api.UploadSubject{Type: "REMUX"},
			want: "2",
		},
		{
			name: "web dl",
			meta: api.UploadSubject{Type: "WEBDL"},
			want: "3",
		},
		{
			name: "encode",
			meta: api.UploadSubject{Type: "ENCODE"},
			want: "15",
		},
		{
			name: "dvd rip",
			meta: api.UploadSubject{Type: "DVDRIP"},
			want: "15",
		},
		{
			name: "unknown",
			meta: api.UploadSubject{Type: "OTHER"},
			want: "0",
		},
	}
	for _, test := range typeTests {
		t.Run(test.name, func(t *testing.T) {
			if got := profile.Site.ResolveTypeID(test.meta); got != test.want {
				t.Fatalf("type ID = %q, want %q", got, test.want)
			}
		})
	}

	resolutionTests := map[string]string{
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
		"1440p": "10",
	}
	for resolution, want := range resolutionTests {
		meta := api.UploadSubject{Release: api.ReleaseInfo{Resolution: resolution}}
		if got := profile.Site.ResolveResolutionID(meta); got != want {
			t.Errorf("resolution %s ID = %q, want %q", resolution, got, want)
		}
	}
}

func TestAdditionalPayload(t *testing.T) {
	data := map[string]string{"internal": "0", "mal": "0"}
	Profile().Site.ApplyAdditionalPayload(trackers.PreparationInput{Meta: api.UploadSubject{
		DiscType:    "BDMV",
		Edition:     "Director's Cut",
		Repack:      "REPACK",
		Region:      "A",
		Distributor: "Example Distributor",
		TVPack:      true,
		Identity:    api.ExternalIdentity{Category: api.CanonicalCategoryTV},
	}}, data)
	for key, want := range map[string]string{
		"edition":     "Director's Cut",
		"release_tag": "REPACK",
		"region":      "A",
		"distributor": "Example Distributor",
		"season_pack": "1",
	} {
		if got := data[key]; got != want {
			t.Errorf("payload %s = %q, want %q", key, got, want)
		}
	}
	if _, exists := data["internal"]; exists {
		t.Fatal("non-internal HUNO payload retained internal field")
	}
}

func TestValidationPolicy(t *testing.T) {
	tests := []struct {
		name     string
		subject  api.TrackerValidationSubject
		wantRule string
	}{
		{
			name:     "blocks web rip",
			subject:  api.TrackerValidationSubject{Type: "WEBRIP"},
			wantRule: "huno_webrip",
		},
		{name: "accepts CRF 22", subject: qualitySubject(api.VideoCRFAssessment{Status: api.VideoCRFStatusPresent, Value: 22}, api.VideoBitrateAssessment{})},
		{
			name:     "blocks CRF above 22",
			subject:  qualitySubject(api.VideoCRFAssessment{Status: api.VideoCRFStatusPresent, Value: 22.1}, api.VideoBitrateAssessment{}),
			wantRule: "huno_video_crf",
		},
		{name: "accepts bitrate fallback", subject: qualitySubject(api.VideoCRFAssessment{Status: api.VideoCRFStatusUnavailable}, api.VideoBitrateAssessment{Status: api.VideoBitrateStatusPresent, BitsPerSecond: 3_000_000})},
		{
			name:     "blocks low bitrate fallback",
			subject:  qualitySubject(api.VideoCRFAssessment{Status: api.VideoCRFStatusUnavailable}, api.VideoBitrateAssessment{Status: api.VideoBitrateStatusPresent, BitsPerSecond: 2_999_999}),
			wantRule: "huno_video_bitrate",
		},
		{name: "allows animation bitrate exception", subject: animationQualitySubject()},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			failures, err := checkRequirements(context.Background(), test.subject, api.NopLogger{})
			if err != nil {
				t.Fatalf("check requirements: %v", err)
			}
			if test.wantRule == "" {
				if len(failures) != 0 {
					t.Fatalf("unexpected failures: %#v", failures)
				}
				return
			}
			if len(failures) != 1 || failures[0].Rule != test.wantRule {
				t.Fatalf("failures = %#v, want rule %q", failures, test.wantRule)
			}
		})
	}
}

func qualitySubject(crf api.VideoCRFAssessment, bitrate api.VideoBitrateAssessment) api.TrackerValidationSubject {
	return api.TrackerValidationSubject{
		Type: "ENCODE",
		Assessments: api.ReleaseAssessments{
			VideoCRF:     crf,
			VideoBitrate: bitrate,
		},
	}
}

func animationQualitySubject() api.TrackerValidationSubject {
	subject := qualitySubject(
		api.VideoCRFAssessment{Status: api.VideoCRFStatusUnavailable},
		api.VideoBitrateAssessment{Status: api.VideoBitrateStatusPresent, BitsPerSecond: 1},
	)
	subject.ProviderMetadata.TMDB = &api.TMDBMetadata{Genres: "Animation, Family"}
	return subject
}
