package module

import (
	"os"
	"strings"
)

func RegX() *LogZ {
	var configPath = os.Getenv("KUBEXDS_CONFIGFILE")
	var keyPath = os.Getenv("KUBEXDS_KEYFILE")
	var certPath = os.Getenv("KUBEXDS_CERTFILE")
	var hideBannerV = os.Getenv("KUBEXDS_HIDEBANNER")

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
