package gtranslate

import (
	"bufio"
	"os"
	"regexp"
	"strings"
)

var stopWords = map[string]bool{
	"a":      true,
	"an":     true,
	"and":    true,
	"as":     true,
	"at":     true,
	"be":     true,
	"but":    true,
	"by":     true,
	"for":    true,
	"if":     true,
	"in":     true,
	"is":     true,
	"it":     true,
	"of":     true,
	"on":     true,
	"or":     true,
	"so":     true,
	"the":    true,
	"to":     true,
	"with":   true,
	"you":    true,
	"your":   true,
	"that":   true,
	"this":   true,
	"from":   true,
	"not":    true,
	"are":    true,
	"have":   true,
	"has":    true,
	"which":  true,
	"they":   true,
	"their":  true,
	"we":     true,
	"us":     true,
	"i":      true,
	"me":     true,
	"my":     true,
	"him":    true,
	"his":    true,
	"her":    true,
	"hers":   true,
	"its":    true,
	"our":    true,
	"ours":   true,
	"yours":  true,
	"thou":   true,
	"thee":   true,
	"thy":    true,
	"do":     true,
	"did":    true,
	"done":   true,
	"does":   true,
	"doing":  true,
	"had":    true,
	"hadst":  true,
	"hath":   true,
	"may":    true,
	"might":  true,
	"must":   true,
	"shall":  true,
	"should": true,
	"will":   true,
	"wilt":   true,
	"would":  true,
}

func SplitIntoWords(text string) []string {
	words := strings.Fields(text)
	var result []string
	for _, word := range words {
		// Remove non-alphabetic characters and convert to lowercase
		word = strings.ToLower(regexp.MustCompile("[^a-zA-Z]+").ReplaceAllString(word, ""))
		if !stopWords[word] && word != "" {
			result = append(result, word)
		}
	}
	return result
}

func loadStopWords(filename string) (map[string]bool, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	stopWords := make(map[string]bool)
	for scanner.Scan() {
		stopWords[scanner.Text()] = true
	}
	return stopWords, scanner.Err()
}

func SplitIntoWordsFile(text string) []string {
	stopWords, err := loadStopWords("english.txt")
	if err != nil {
		// handle error
	}
	words := strings.Fields(text)
	var result []string
	wordMap := make(map[string]bool)
	for _, word := range words {
		// Remove non-alphabetic characters and convert to lowercase
		word = strings.ToLower(regexp.MustCompile("[^a-zA-Z]+").ReplaceAllString(word, ""))
		wordMap[word] = false
	}

	for _, word := range words {
		// Remove non-alphabetic characters and convert to lowercase
		word = strings.ToLower(regexp.MustCompile("[^a-zA-Z]+").ReplaceAllString(word, ""))
		if !wordMap[word] && !stopWords[word] && word != "" {
			result = append(result, word)
			wordMap[word] = true
		}
	}
	return result
}

func SplitIntoSentences(text string) []string {
	var sentences []string

	// Split text into individual words
	words := strings.Fields(text)

	// Define punctuation marks that typically end a sentence
	punctuation := ".?!;,"

	// Build a slice of sentences
	var sentence strings.Builder
	for _, word := range words {
		sentence.WriteString(word)
		sentence.WriteString(" ")

		// If we encounter a punctuation mark that ends a sentence, add the current sentence to the slice
		if strings.ContainsAny(string(word[len(word)-1]), punctuation) {
			sentenceStr := strings.TrimSpace(sentence.String())
			if len(sentenceStr) > 0 {
				sentences = append(sentences, sentenceStr)
			}
			sentence.Reset()
		}
	}

	// Add any remaining text to the slice
	if sentence.Len() > 0 {
		sentences = append(sentences, strings.TrimSpace(sentence.String()))
	}

	return sentences
}
