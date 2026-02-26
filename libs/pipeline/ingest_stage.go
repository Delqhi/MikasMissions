package pipeline

import (
	contractsevents "github.com/delqhi/mikasmissions/platform/libs/contracts-events"
)

func BuildTranscodeRequest(input contractsevents.MediaUploadedV1) (contractsevents.MediaTranscodeRequestedV1, error) {
	if err := input.Validate(); err != nil {
		return contractsevents.MediaTranscodeRequestedV1{}, err
	}
	return contractsevents.MediaTranscodeRequestedV1{
		AssetID:   input.AssetID,
		SourceURL: input.SourceURL,
		TraceID:   input.TraceID,
	}, nil
}
