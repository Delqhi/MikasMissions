package contractsevents

import "errors"

type MediaTranscodeRequestedV1 struct {
	AssetID   string `json:"asset_id"`
	SourceURL string `json:"source_url"`
	TraceID   string `json:"trace_id"`
}

func (e MediaTranscodeRequestedV1) Validate() error {
	if e.AssetID == "" || e.SourceURL == "" || e.TraceID == "" {
		return errors.New("media.transcode.requested.v1 has missing required fields")
	}
	return nil
}
