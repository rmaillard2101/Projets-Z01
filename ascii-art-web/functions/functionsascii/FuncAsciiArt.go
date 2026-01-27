package artweb

import (
	art "ascii-art-web/functions/functionsascii/functions"
	"fmt"
	"strings"
)

func AsciiArt(sentence string, bannerfile string) [][]string {
	var rettab [][]string
	sentence = strings.ReplaceAll(sentence, `\n`, "\n")
	artsentence := []string{}
	for i := 0; i < len(sentence); i++ {
		if sentence[i] != 10 {
			convertedletter := art.ConvertitEnMatriceSansBalisage(sentence[i], bannerfile)
			artsentence = art.AddLetter(artsentence, convertedletter)

		} else {
			if len(artsentence) != 0 {
				rettab = append(rettab, artsentence)
				artsentence = []string{}
			} else {
				fmt.Println("")
			}

		}
	}
	rettab = append(rettab, artsentence)
	return rettab
}
