package core

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

type ImportantExtractor struct {
	UniModelProb UnigramProabilityCollections
}

func NewImportantExtractor(unictx UniGramModel) ImportantExtractor {
	return ImportantExtractor{
		UniModelProb: unictx.GetProabilityWeight(),
	}
}

func (ie *ImportantExtractor) Extract(rawinput string) tokenTupleGroups {
	input := stringFilter(rawinput)

	extractUnimodel := NewUniGramModel()
	extractUnimodel.Update(strings.Join(RemoveStopwords(strings.Split(input, " ")), " "))
	targerUniData := extractUnimodel.GetProabilityWeight()

	candidates := make(tokenTupleGroups, 0)

	trimedstrings := make([]string, 0)
	for _, strings := range extractUnimodel.gramlize(input) {
		trimedstrings = append(trimedstrings, strings[0])
	}
	trimedstrings = RemoveStopwords(trimedstrings)
	fmt.Println(trimedstrings)

	var totalLength int = 0
	for _, word := range trimedstrings {
		totalLength += len(word)
	}

	var lengthAvg float64 = (float64(totalLength) / float64(len(trimedstrings)))

	for formerdata, _ := range targerUniData {
		if formerdata == ENDTOKEN {
			break
		}

		var prob float64
		for _, proablity := range ie.UniModelProb[formerdata] {
			prob += proablity
		}

		prob += float64((freqInArray(formerdata, trimedstrings) / len(trimedstrings)))

		if float64(len(formerdata)) > lengthAvg {
			prob += float64(len(formerdata)) - lengthAvg
		}

		token := tokenTuple{
			Token:      formerdata,
			Proability: prob,
		}
		candidates = append(candidates, token)
	}

	sort.Sort(sort.Reverse(candidates))
	fmt.Println(candidates)

	return candidates
}

func stringFilter(input string) string {
	//'https?:\/\/.*?[\s+]'gm
	removeLink := regexp.MustCompile(`'https?:\/\/.*?[\s+]'gm`)
	removeMention := regexp.MustCompile(`\S*@\S*\S*@\S*\s?`)

	linkClear := removeLink.ReplaceAllString(input, "")
	mentionClear := removeMention.ReplaceAllString(linkClear, "")

	return mentionClear
}
