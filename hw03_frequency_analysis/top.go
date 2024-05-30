package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

type wordFrequency struct {
	word  string
	count int
}

func SplitText(text string) []string {
	return strings.Fields(text)
}

func CountWordFrequency(words []string) map[string]int {
	wordCount := make(map[string]int)
	for _, word := range words {
		wordCount[word]++
	}
	return wordCount
}

func SortWordFrequencies(frequencies []wordFrequency) {
	sort.Slice(frequencies, func(i, j int) bool {
		if frequencies[i].count == frequencies[j].count {
			return frequencies[i].word < frequencies[j].word
		}
		return frequencies[i].count > frequencies[j].count
	})
}

func GetTopTenWords(frequencies []wordFrequency) []string {
	// Sort the frequencies slice based on frequency in descending order
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

func Top10(text string) []string {
	words := SplitText(text)
	wordCount := CountWordFrequency(words)

	frequencies := make([]wordFrequency, 0, len(wordCount))
	for word, count := range wordCount {
		frequencies = append(frequencies, wordFrequency{word, count})
	}

	SortWordFrequencies(frequencies)

	return GetTopTenWords(frequencies)
}
