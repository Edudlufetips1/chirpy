package main

import "strings"

func filterProfanity(chirp string) string {
	isProfane := func(word string) bool {
		profanityList := []string{"kerfuffle", "sharbert", "fornax"}
		lowerWord := strings.ToLower(word)
		for _, profaneWord := range profanityList {
			if lowerWord == profaneWord {
				return true
			}
		}
		return false
	}
	text := strings.Split(chirp, " ")
	for i, word := range text {
		if isProfane(word) {
			text[i] = "****"
		}
	}
	return strings.Join(text, " ")
}
