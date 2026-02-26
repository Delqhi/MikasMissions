package pipeline

import (
	"testing"

	contractsevents "github.com/delqhi/mikasmissions/platform/libs/contracts-events"
)

func TestPipelineChainEndToEnd(t *testing.T) {
	uploaded := contractsevents.MediaUploadedV1{
		AssetID:   "asset-123",
		SourceURL: "https://cdn.local/asset-123.mp4",
		Uploader:  "u-1",
		TraceID:   "trace-1",
	}
	transcodeReq, err := BuildTranscodeRequest(uploaded)
	if err != nil {
		t.Fatalf("BuildTranscodeRequest: %v", err)
	}
	transcoded, err := BuildTranscodedMedia(transcodeReq)
	if err != nil {
		t.Fatalf("BuildTranscodedMedia: %v", err)
	}
	reviewed, approved, err := BuildPolicyOutputs(transcoded)
	if err != nil {
		t.Fatalf("BuildPolicyOutputs: %v", err)
	}
	published, err := BuildEpisodePublished(approved)
	if err != nil {
		t.Fatalf("BuildEpisodePublished: %v", err)
	}
	if reviewed.PolicyResult != "approved" {
		t.Fatalf("expected policy result approved, got %s", reviewed.PolicyResult)
	}
	if published.EpisodeID == "" {
		t.Fatalf("expected non-empty episode id")
	}
	if len(published.LearningTags) == 0 {
		t.Fatalf("expected learning tags")
	}
}
