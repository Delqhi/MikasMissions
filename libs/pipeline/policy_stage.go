package pipeline

import (
	contractsevents "github.com/delqhi/mikasmissions/platform/libs/contracts-events"
)

func BuildPolicyOutputs(input contractsevents.MediaTranscodedV1) (contractsevents.MediaReviewedV1, contractsevents.MediaApprovedV1, error) {
	if err := input.Validate(); err != nil {
		return contractsevents.MediaReviewedV1{}, contractsevents.MediaApprovedV1{}, err
	}
	reviewed := contractsevents.MediaReviewedV1{
		AssetID:      input.AssetID,
		PolicyResult: "approved",
		AgeBand:      "6-11",
	}
	approved := contractsevents.MediaApprovedV1{
		AssetID:      input.AssetID,
		AgeBand:      "6-11",
		LearningTags: []string{"farben", "teamwork"},
	}
	return reviewed, approved, nil
}
