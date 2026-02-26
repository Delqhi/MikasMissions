package contractsevents

import "errors"

type MediaRendition struct {
	Profile string `json:"profile"`
	URL     string `json:"url"`
}

type MediaTranscodedV1 struct {
	AssetID    string           `json:"asset_id"`
	Renditions []MediaRendition `json:"renditions"`
	DurationMS int64            `json:"duration_ms"`
}

func (e MediaTranscodedV1) Validate() error {
	if e.AssetID == "" || len(e.Renditions) == 0 || e.DurationMS <= 0 {
		return errors.New("media.transcoded.v1 has invalid payload")
	}
	return nil
}
