package cli

import (
	"os"
	"strings"
)

func GetDescriptions(descriptionArg []string, hideBanner bool) map[string]string {
	var description, banner string

	if strings.Contains(strings.Join(os.Args[0:], ""), "-h") {
		description = descriptionArg[0]
	} else {
		if len(descriptionArg) > 1 {
			description = descriptionArg[1]
		} else {
			description = descriptionArg[0]
		}
	}

	if !hideBanner {
		banner = `
 __                                  
|  \                                 
| ▓▓       ______   ______  ________ 
| ▓▓      /      \ /      \|        \
| ▓▓     |  ▓▓▓▓▓▓\  ▓▓▓▓▓▓\\▓▓▓▓▓▓▓▓
| ▓▓     | ▓▓  | ▓▓ ▓▓  | ▓▓ /    ▓▓ 
| ▓▓_____| ▓▓__/ ▓▓ ▓▓__| ▓▓/  ▓▓▓▓_ 
| ▓▓     \\▓▓    ▓▓\▓▓    ▓▓  ▓▓    \
 \▓▓▓▓▓▓▓▓ \▓▓▓▓▓▓ _\▓▓▓▓▓▓▓\▓▓▓▓▓▓▓▓
                  |  \__| ▓▓         
                   \▓▓    ▓▓         
                    \▓▓▓▓▓▓          `
	} else {
		banner = ""
	}
	return map[string]string{"banner": banner, "description": description}
}
