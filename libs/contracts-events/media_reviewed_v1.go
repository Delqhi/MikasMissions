package contractsevents

import "errors"

type MediaReviewedV1 struct {
	AssetID      string `json:"asset_id"`
	PolicyResult string `json:"policy_result"`
	AgeBand      string `json:"age_band"`
}

func (e MediaReviewedV1) Validate() error {
	if e.AssetID == "" || e.PolicyResult == "" || e.AgeBand == "" {
		return errors.New("media.reviewed.v1 has missing required fields")
	}
	return nil
}
