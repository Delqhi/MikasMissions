package contractsevents

import "errors"

type MediaApprovedV1 struct {
	AssetID      string   `json:"asset_id"`
	AgeBand      string   `json:"age_band"`
	LearningTags []string `json:"learning_tags"`
}

func (e MediaApprovedV1) Validate() error {
	if e.AssetID == "" || e.AgeBand == "" {
		return errors.New("media.approved.v1 has missing required fields")
	}
	return nil
}
