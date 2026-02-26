package contractsevents

import (
	"encoding/json"
	"testing"
)

func TestVideoAssetReadyV1Contract(t *testing.T) {
	raw := []byte(`{"run_id":"run1","asset_id":"asset1","source_url":"https://cdn.example/asset1.mp4","duration_ms":120000,"content_suitability":"core","age_band":"6-11","uploader_id":"admin1","ready_at":"2026-03-11T09:00:00Z"}`)
	var event VideoAssetReadyV1
	if err := json.Unmarshal(raw, &event); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if err := event.Validate(); err != nil {
		t.Fatalf("validate: %v", err)
	}
}

func TestVideoAssetReadyV1ContractRejectsDuration(t *testing.T) {
	event := VideoAssetReadyV1{RunID: "run1", AssetID: "asset1", SourceURL: "https://cdn.example/a.mp4", ContentSuitability: "core", AgeBand: "6-11", ReadyAt: "2026-03-11T09:00:00Z"}
	if err := event.Validate(); err == nil {
		t.Fatalf("expected validation error")
	}
}
