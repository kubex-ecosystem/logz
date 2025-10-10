package module

import (
	"os"
)

func RegX() *Logz {
	return &Logz{
		hideBanner: os.Getenv("GOBE_HIDE_BANNER") == "true",
	}
}
