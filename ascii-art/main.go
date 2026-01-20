package main

import (
	art "ascii-art/functions"
	"fmt"
	"os"
	"strings"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Need a string as argument")
		return
	} else if len(os.Args) > 2 {
		fmt.Println("Too many arguments : need only one string")
		return
	}
	sentence := strings.ReplaceAll(os.Args[1], `\n`, "\n")
	artsentence := []string{}
	for i := 0; i < len(sentence); i++ {
		if sentence[i] != 10 {
			convertedletter := art.ConvertitEnMatriceSansBalisage(sentence[i], "standard.txt")
			artsentence = art.AddLetter(artsentence, convertedletter)

		} else {
			if len(artsentence) != 0 {
				printArt(artsentence)
				artsentence = []string{}
			} else {
				fmt.Println("")
			}

		}
	}
	printArt(artsentence)
}

func printArt(artsentence []string) {
	for _, line := range artsentence {
		fmt.Println(line)
	}
}
