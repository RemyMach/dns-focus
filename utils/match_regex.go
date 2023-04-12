package utils

import (
	"fmt"
	"regexp"
)

func MatchRegex(pattern, text string) bool {
	re, err := regexp.Compile(pattern)
	if err != nil {
		fmt.Printf("Erreur lors de la compilation de l'expression régulière: %v\n", err)
		return false
	}

	return re.MatchString(text)
}