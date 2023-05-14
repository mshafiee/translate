package gtranslate

import (
	"fmt"
	"golang.org/x/text/language"
	"time"
)

var GoogleHost = "google.com"

// TranslationParams is a util struct to pass as parameter to indicate how to translate
type TranslationParams struct {
	From       string
	To         string
	Tries      int
	Delay      time.Duration
	GoogleHost string
}

// Translate translate a text using native tags offer by go language
func Translate(text string, from language.Tag, to language.Tag, googleHost ...string) ([]string, error) {
	if len(googleHost) != 0 && googleHost[0] != "" {
		GoogleHost = googleHost[0]
	}
	translated, err := translate(text, from.String(), to.String(), false, 2, 0)
	if err != nil {
		return nil, err
	}

	return translated, nil
}

// TranslateWithParams translate a text with simple params as string
func TranslateWithParams(text string, params TranslationParams) ([]string, error) {
	if params.GoogleHost == "" {
		GoogleHost = "google.com"
	} else {
		GoogleHost = params.GoogleHost
	}
	translated, err := translate(text, params.From, params.To, true, params.Tries, params.Delay)
	if err != nil {
		return nil, err
	}
	return translated, nil
}

// VocabularyWithParams translate vocabulary of text
func VocabularyWithParams(text string, params TranslationParams) ([]string, error) {
	var vocabularyMeaning []string

	if params.GoogleHost == "" {
		GoogleHost = "google.com"
	} else {
		GoogleHost = params.GoogleHost
	}

	words := SplitIntoWordsFile(text)

	for _, w := range words {
		translated, err := translate(w, params.From, params.To, true, params.Tries, params.Delay)
		if err != nil {
			return nil, err
		}

		vocab := fmt.Sprintf("%s: ", w)
		for i, t := range translated {
			if i != len(translated)-1 {
				vocab += fmt.Sprintf("%s,", t)
			} else {
				vocab += fmt.Sprintf("%s.", t)
			}
		}
		vocabularyMeaning = append(vocabularyMeaning, vocab)
	}

	return vocabularyMeaning, nil
}

// SentenceWithParams translate sentences of text
func SentenceWithParams(text string, params TranslationParams) ([]string, error) {
	var sentenceMeaning []string

	if params.GoogleHost == "" {
		GoogleHost = "google.com"
	} else {
		GoogleHost = params.GoogleHost
	}

	sentences := SplitIntoSentences(text)

	for _, s := range sentences {
		translated, err := translate(s, params.From, params.To, true, params.Tries, params.Delay)
		if err != nil {
			return nil, err
		}

		// Create a map to keep track of which strings we have seen
		seen := make(map[string]bool)

		sentence := fmt.Sprintf("%s: ", s)
		for _, t := range translated {
			if !seen[t] {
				sentence += fmt.Sprintf("%s ", t)
				seen[t] = true
			}
		}
		sentenceMeaning = append(sentenceMeaning, sentence)
	}

	return sentenceMeaning, nil
}
