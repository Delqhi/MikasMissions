package contractsevents

import "errors"

type MediaUploadedV1 struct {
	AssetID   string `json:"asset_id"`
	SourceURL string `json:"source_url"`
	Uploader  string `json:"uploader_id"`
	TraceID   string `json:"trace_id"`
}

func (e MediaUploadedV1) Validate() error {
	if e.AssetID == "" || e.SourceURL == "" || e.Uploader == "" || e.TraceID == "" {
		return errors.New("media.uploaded.v1 has missing required fields")
	}
	return nil
}
