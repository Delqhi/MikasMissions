package contractsevents

import (
	"encoding/json"
	"testing"
)

func TestMediaUploadedV1Contract(t *testing.T) {
	raw := []byte(`{"asset_id":"a1","source_url":"https://cdn/x.mp4","uploader_id":"u1","trace_id":"tr1"}`)
	var event MediaUploadedV1
	if err := json.Unmarshal(raw, &event); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if err := event.Validate(); err != nil {
		t.Fatalf("validate: %v", err)
	}
}
