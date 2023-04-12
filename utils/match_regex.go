package utils

import (
	"fmt"
	"regexp"
)

func MatchRegex(pattern, text string) bool {
	// Compile l'expression régulière
	re, err := regexp.Compile(pattern)
	if err != nil {
		fmt.Printf("Erreur lors de la compilation de l'expression régulière: %v\n", err)
		return false
	}

	// Vérifie si le texte correspond à l'expression régulière
	return re.MatchString(text)
}