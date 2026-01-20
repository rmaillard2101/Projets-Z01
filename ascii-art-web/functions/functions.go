package art2

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func RecoverLetter(letter rune, nameFile string) []string {
	file, err := os.Open(nameFile)
	if err != nil {
		return nil
	}
	defer file.Close()
	var linesTable []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		linesTable = append(linesTable, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil
	}
	letterArt := make([]string, 8)
	index := (int(letter) - 32) * 9
	for i := 0; i < 8; i++ {
		letterArt[i] = linesTable[index+i+1]
	}

	return letterArt
}

func AddLetter(sentence []string, letter []string) []string {
	if len(sentence) == 0 {
		return letter
	}
	for i := range letter {
		sentence[i] += letter[i]
	}
	return sentence
}

// Version modifiée de Ascii-Art pour renvoyer un string
func TransformPhrase(sentence string, banner string) (string, error) {
	if len(sentence) > 1000 {
		return "Sentence too long, may crash the server.", fmt.Errorf("phrase trop longue")
	}
	sentence = strings.ReplaceAll(sentence, `\n`, "\n")
	sentenceArt := []string{}
	var result strings.Builder

	nameFile := "standard.txt"
	if banner != "" {
		nameFile = banner + ".txt"
		if RecoverLetter('a', nameFile) == nil {
			return "Could not read file : verify spelling and existence of the wanted file\n", nil
		}
	}

	for _, letter := range sentence {
		if letter != '\n' {
			letterArt := RecoverLetter(letter, nameFile)
			sentenceArt = AddLetter(sentenceArt, letterArt)
		} else {
			if len(sentenceArt) != 0 {
				for _, line := range sentenceArt {
					result.WriteString(line + "\n")
				}
				sentenceArt = []string{}
			} else {
				result.WriteString("\n")
			}
		}
	}
	if len(sentenceArt) != 0 {
		for _, line := range sentenceArt {
			result.WriteString(line + "\n")
		}
	}

	return result.String(), nil
}
