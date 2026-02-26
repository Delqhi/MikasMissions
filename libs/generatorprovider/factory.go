package generatorprovider

import "fmt"

func NewProvider(profile ModelProfile) (Provider, error) {
	switch profile.Provider {
	case "nvidia_nim":
		return NewNIMProvider(profile), nil
	default:
		return nil, fmt.Errorf("unsupported provider: %s", profile.Provider)
	}
}
