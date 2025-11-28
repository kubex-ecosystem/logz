package module

import (
	"os"
	"strings"
)

func RegX() *LogZ {
	var configPath = os.Getenv("CANALIZEDS_CONFIGFILE")
	var keyPath = os.Getenv("CANALIZEDS_KEYFILE")
	var certPath = os.Getenv("CANALIZEDS_CERTFILE")
	var hideBannerV = os.Getenv("CANALIZEDS_HIDEBANNER")

	return &LogZ{
		configPath: configPath,
		keyPath:    keyPath,
		certPath:   certPath,
		hideBanner: (strings.ToLower(hideBannerV) == "true" ||
			strings.ToLower(hideBannerV) == "1" ||
			strings.ToLower(hideBannerV) == "yes" ||
			strings.ToLower(hideBannerV) == "y"),
	}
}
