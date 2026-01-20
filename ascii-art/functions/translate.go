package art

func AddLetter(sentence []string, letter []string) []string {
	if len(sentence) == 0 {
		return letter
	}
	for i := range letter {
		sentence[i] += letter[i]
	}
	return sentence
}
