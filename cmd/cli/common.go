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
 ___           ________      ________      ________     
|\  \         |\   __  \    |\   ____\    |\_____  \    
\ \  \        \ \  \|\  \   \ \  \___|     \|___/  /|   
 \ \  \        \ \  \\\  \   \ \  \  ___       /  / /   
  \ \  \____    \ \  \\\  \   \ \  \|\  \     /  /_/__  
   \ \_______\   \ \_______\   \ \_______\   |\________\
    \|_______|    \|_______|    \|_______|    \|_______|
`
	} else {
		banner = ""
	}
	return map[string]string{"banner": banner, "description": description}
}
