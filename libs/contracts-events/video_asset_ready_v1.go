package contractsevents

import "errors"

type VideoAssetReadyV1 struct {
	RunID              string `json:"run_id"`
	AssetID            string `json:"asset_id"`
	SourceURL          string `json:"source_url"`
	DurationMS         int64  `json:"duration_ms"`
	ContentSuitability string `json:"content_suitability"`
	AgeBand            string `json:"age_band"`
	UploaderID         string `json:"uploader_id"`
	ReadyAt            string `json:"ready_at"`
}

func (e VideoAssetReadyV1) Validate() error {
	if e.RunID == "" || e.AssetID == "" || e.SourceURL == "" || e.ReadyAt == "" {
		return errors.New("video.asset.ready.v1 has missing required fields")
	}
	if e.ContentSuitability == "" || e.AgeBand == "" {
		return errors.New("video.asset.ready.v1 has missing content fields")
	}
	if e.DurationMS <= 0 {
		return errors.New("video.asset.ready.v1 requires positive duration_ms")
	}
	return nil
}
