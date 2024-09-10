package core

import (
	"encoding/json"
	"os"
	"strings"
)

type stopwordsData struct {
	Stopwords []string `json:"stopwords"`
}

// temp
var Stopwords []string = make([]string, 0)

func init() {
	loadStopwords()
}

func loadStopwords() {
	data, ferr := os.ReadFile("./stop.json")
	if ferr != nil {
		panic(ferr)
	}

	words := *new(stopwordsData)
	jerr := json.Unmarshal(data, &words)
	if jerr != nil {
		panic(jerr)
	}

	Stopwords = append(Stopwords, words.Stopwords...)
}

func RemoveStopwords(input []string) []string {
	filtered := make([]string, 0)

	for _, words := range input {
		for _, stop := range Stopwords {

			if strings.HasSuffix(words, stop) {
				cutedword, _ := strings.CutSuffix(words, stop)
				filtered = append(filtered, cutedword)
			}
		}
	}

	return filtered
}
