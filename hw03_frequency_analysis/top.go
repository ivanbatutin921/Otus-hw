package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

func Top10(text string) []string {
	words := strings.Fields(text)
	wordCount := make(map[string]int)

	for _, word := range words {
		wordCount[word]++
	}

	type wordFrequency struct {
		word  string
		count int
	}

	frequencies := make([]wordFrequency, 0, len(wordCount))

	for word, count := range wordCount {
		frequencies = append(frequencies, wordFrequency{word, count})
	}

	sort.Slice(frequencies, func(i, j int) bool {
		if frequencies[i].count == frequencies[j].count {
			return frequencies[i].word < frequencies[j].word
		}
		return frequencies[i].count > frequencies[j].count
	})

	var topTen []string
	for i := 0; i < len(frequencies) && i < 10; i++ {
		topTen = append(topTen, frequencies[i].word)
	}
	return topTen
}
