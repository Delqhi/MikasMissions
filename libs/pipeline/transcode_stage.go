package pipeline

import (
	"fmt"

	contractsevents "github.com/delqhi/mikasmissions/platform/libs/contracts-events"
)

func BuildTranscodedMedia(input contractsevents.MediaTranscodeRequestedV1) (contractsevents.MediaTranscodedV1, error) {
	if err := input.Validate(); err != nil {
		return contractsevents.MediaTranscodedV1{}, err
	}
	return contractsevents.MediaTranscodedV1{
		AssetID: input.AssetID,
		Renditions: []contractsevents.MediaRendition{
			{Profile: "1080p", URL: fmt.Sprintf("https://cdn.example.local/%s/1080p.m3u8", input.AssetID)},
			{Profile: "720p", URL: fmt.Sprintf("https://cdn.example.local/%s/720p.m3u8", input.AssetID)},
		},
		DurationMS: 660000,
	}, nil
}
