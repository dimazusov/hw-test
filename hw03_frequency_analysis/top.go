package hw03_frequency_analysis //nolint:golint,stylecheck

import (
	"errors"
	"log"
	"regexp"
	"sort"
	"strings"
)

var ErrWordNotFound = errors.New("Word not found")

func Top10(text string) []string {
	if text == "" {
		return []string{}
	}

	wordsMap := splitTextToWords(&text)

	maxLenMostFreqWords := 12
	resMap := make(map[int][]string)

	if len(wordsMap) < maxLenMostFreqWords {
		maxLenMostFreqWords = len(wordsMap)
	}


	for i:=0; i<maxLenMostFreqWords; i++ {
		count, bestFreqWord, err := getBestFreqWord(wordsMap)
		if err != nil {
			log.Fatalln(err)
		}

		resMap[count] = append(resMap[count], bestFreqWord)
		delete(wordsMap, bestFreqWord)
	}

	wordCounts := []int{}
	for i, words := range resMap {
		sort.Strings(words)
		wordCounts = append(wordCounts, i)
	}
	sort.Slice(wordCounts, func (i, j int) bool {return i < j})

	mostFreqWords := []string{}
	for _, wordCount := range wordCounts {
		mostFreqWords = append(mostFreqWords, resMap[wordCount]...)
	}

	return mostFreqWords
}

func getBestFreqWord(hm map[string]int) (count int, word string, err error) {
	maxCount := 0
	maxCountWord := ""
	for word, count := range hm {
		if count > maxCount {
			maxCount = count
			maxCountWord = word
		}
	}

	if maxCountWord == "" {
		return 0, "", ErrWordNotFound
	}

	return maxCount, maxCountWord, nil
}

func splitTextToWords(text *string) map[string]int {
	var re = regexp.MustCompile(`[\n\s]+`)
	textWithoutLineBreakAndTabs := re.ReplaceAllString(*text, " ")

	wordsMap := make(map[string]int)
	re = regexp.MustCompile(`[,\.\";(!-):]|^\-$`)

	for _, word := range strings.Split(textWithoutLineBreakAndTabs, " ") {
		replacableWord := re.ReplaceAllString(strings.ToLower(word), "")

		if len(replacableWord) != 0 {
			wordsMap[replacableWord]++
		}
	}

	return wordsMap
}