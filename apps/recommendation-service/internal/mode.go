package internal

import contractsapi "github.com/delqhi/mikasmissions/platform/libs/contracts-api"

func normalizeMode(mode string) string {
	switch mode {
	case "core", "teen":
		return mode
	default:
		return "early"
	}
}

func ageBandForMode(mode string) string {
	switch mode {
	case "core":
		return contractsapi.AgeBandCore
	case "teen":
		return contractsapi.AgeBandTeen
	default:
		return contractsapi.AgeBandEarly
	}
}

func actionsForMode(mode string) []string {
	switch mode {
	case "core":
		return []string{"Resume", "Mission", "Explore", "Progress", "Library"}
	case "teen":
		return []string{"Watch", "Watchlist", "Explore", "Learn", "Report"}
	default:
		return []string{"Start", "Mission", "Favorites", "Pause"}
	}
}
