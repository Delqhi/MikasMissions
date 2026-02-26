package contractsevents

import (
	"encoding/json"
	"testing"
)

func TestMediaTranscodedV1Contract(t *testing.T) {
	raw := []byte(`{"asset_id":"a1","renditions":[{"profile":"720p","url":"https://cdn/720.m3u8"}],"duration_ms":120000}`)
	var event MediaTranscodedV1
	if err := json.Unmarshal(raw, &event); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if err := event.Validate(); err != nil {
		t.Fatalf("validate: %v", err)
	}
}
