package hw03_frequency_analysis //nolint:golint,stylecheck

import (
	"errors"
	"regexp"
	"sort"
	"strings"
)

var ErrWordNotFound = errors.New("word not found")

var RegexEndLine = regexp.MustCompile(`[\n\s]+`)
var RegexCleanSymbol = regexp.MustCompile(`[,\.\";(!-):]|^\-$`)

type MostReqWord struct{
	Word string
	Count int
}

func Top10(text string) []string {
	if text == "" {
		return nil
	}

	words := splitTextToWords(text)

	sort.SliceStable(words, func(i, j int) bool {
		return words[i].Count > words[j].Count
	})

	return getMostRequreciesWords(words, 10)
}

func splitTextToWords(text string) (mostRequrencyWords []MostReqWord) {
	words := strings.Split(RegexEndLine.ReplaceAllString(text, " "), " ")

	for _, word := range words {
		replacableWord := RegexCleanSymbol.ReplaceAllString(strings.ToLower(word), "")

		if len(replacableWord) == 0 {
			continue
		}

		index, ok := findIndexByWord(mostRequrencyWords, replacableWord)
		if ok {
			mostRequrencyWords[index].Count++
		} else {
			mostRequrencyWords = append(mostRequrencyWords, MostReqWord{Word: replacableWord, Count: 1})
		}
	}

	return mostRequrencyWords
}

func findIndexByWord (mostRequrencyWords []MostReqWord, word string) (int, bool) {
	for i, mostReqWord := range mostRequrencyWords {
		if mostReqWord.Word == word {
			return i, true
		}
	}

	return 0, false
}

func getMostRequreciesWords(reqWords []MostReqWord, maxCountMostRequreciesWords int) (mostReqWords []string) {
	for i, word := range reqWords {
		mostReqWords = append(mostReqWords, word.Word)

		if i >= maxCountMostRequreciesWords && !isNextWordHasSameCount(reqWords, i) {
			break
		}
	}

	return mostReqWords
}

func isNextWordHasSameCount(reqWords []MostReqWord, index int) bool {
	if len(reqWords) - 1 == index {
		return false
	}

	if reqWords[index].Count != reqWords[index+1].Count {
		return false
	}

	return true
}