package core

import (
	"fmt"
	"strings"
)

// type BigramContextCollections map[string]map[string]string

//type TokenRelationMaps map[string]map[string]float64

const ENDTOKEN = "[<EOT>]"

type UniGramModel struct {
	//Context   BigramContextCollections
	tokenProability    UnigramProabilityCollections
	tokenAppearContext uniGramContextTokenAppears

	TokenProabilityWeight UnigramProabilityCollections `json:"weight"`

	Pretrained bool
}

func NewUniGramModel() UniGramModel {
	return UniGramModel{
		tokenProability:       make(UnigramProabilityCollections),
		TokenProabilityWeight: make(UnigramProabilityCollections),
		tokenAppearContext:    make(uniGramContextTokenAppears),
		Pretrained:            false,
	}
}

func (nm *UniGramModel) gramlize(str string) [][]string {
	//will be replaced by improved lexer.
	str = strings.ToLower(str)
	splited := strings.Split(str, " ")

	var bigram_collections = make([][]string, 0)

	var offset int64 = 0
	var length = len(splited)

	for offset < int64(length) {
		bigram_tuple := []string{}

		if (offset + 1) >= int64(length) {
			bigram_tuple = append(bigram_tuple, splited[offset])
			bigram_tuple = append(bigram_tuple, ENDTOKEN)
		} else {
			bigram_tuple = append(bigram_tuple, splited[offset])
			bigram_tuple = append(bigram_tuple, splited[offset+1])
		}

		bigram_collections = append(bigram_collections, bigram_tuple)
		offset++
	}

	return bigram_collections
}

func (nm *UniGramModel) Update(str string) {
	//Will be replaced
	tokens := nm.gramlize(str)

	for _, bigram_tuples := range tokens {
		current_token := bigram_tuples[0]
		next_token := bigram_tuples[1]

		if _, exist := nm.tokenProability[current_token]; !exist {
			nm.tokenProability[current_token] = make(map[string]float64)
			nm.tokenProability[current_token][next_token] = 1.0
		} else {
			nm.tokenProability[current_token][next_token] += 1.0
		}

		if _, exist := nm.tokenAppearContext[current_token]; !exist {
			nm.tokenAppearContext[current_token] = make([]string, 0)
			nm.tokenAppearContext[current_token] = append(nm.tokenAppearContext[current_token], next_token)
		} else {
			nm.tokenAppearContext[current_token] = append(nm.tokenAppearContext[current_token], next_token)
		}
	}

	//nm.calucateFullProability()
}

func (nm *UniGramModel) GetProabilityWeight() UnigramProabilityCollections {
	if nm.Pretrained {
		return nm.TokenProabilityWeight
	}

	nm.calucateFullProability()
	return nm.TokenProabilityWeight
}

func (nm *UniGramModel) GetSize() int {
	return len(nm.tokenProability)
}

func (nm *UniGramModel) calculateTokenWeight(prevtoken string, nexttoken string) float64 {
	context_total := float64(len(nm.tokenAppearContext[prevtoken]))
	proability_total := nm.tokenProability[prevtoken][nexttoken]

	return proability_total / context_total
}

func (nm *UniGramModel) calucateFullProability() {
	fmt.Println("Calculating Weight... \nSize: ", len(nm.tokenProability))
	for keyprevtoken, internal_tokenmap := range nm.tokenProability {
		for keynexttoken := range internal_tokenmap {
			nm.TokenProabilityWeight[keyprevtoken] = make(map[string]float64)
			nm.TokenProabilityWeight[keyprevtoken][keynexttoken] = nm.calculateTokenWeight(keyprevtoken, keynexttoken)
		}
	}
}

func biTuple(first string, second string) BiGramTokenTuple {
	return BiGramTokenTuple{FirstToken: first, NextToken: second}
}

type BiGramModel struct {
	tokenProability    BigramProabilityCollections
	tokenAppearContext biGramContextTokenAppears

	TokenProabilityWeight BigramProabilityCollections

	Pretrained bool
}

func NewBiGramModel() BiGramModel {
	return BiGramModel{
		tokenProability:    make(BigramProabilityCollections),
		tokenAppearContext: make(biGramContextTokenAppears),

		TokenProabilityWeight: make(BigramProabilityCollections),
		Pretrained:            false,
	}
}

func (nm *BiGramModel) gramlize(str string) []BiGramTokenTuple {
	//will be replaced by improved lexer.
	lowered := strings.ToLower(str)
	splited := strings.Split(lowered, " ")

	var bigram_collections = make([]BiGramTokenTuple, 0)

	var offset int64 = 0
	var length = len(splited)

	for offset < int64(length) {
		var bigram_tuple BiGramTokenTuple

		if (offset + 1) >= int64(length) {
			bigram_tuple = biTuple(splited[offset], ENDTOKEN)
		} else {
			bigram_tuple = biTuple(splited[offset], splited[offset+1])
		}

		bigram_collections = append(bigram_collections, bigram_tuple)
		offset++
	}

	return bigram_collections
}

func (nm *BiGramModel) Update(str string) {
	//Will be replaced
	tokens := nm.gramlize(str)

	for i, bigram_tuples := range tokens {
		var following_token string = ""
		if (i + 1) >= len(tokens) {
			following_token = bigram_tuples.NextToken
		} else {
			following_token = tokens[i+1].FirstToken
		}

		if _, exist := nm.tokenProability[bigram_tuples]; !exist {
			nm.tokenProability[bigram_tuples] = make(map[string]float64)
			nm.tokenProability[bigram_tuples][following_token] = 1.0
		} else {
			nm.tokenProability[bigram_tuples][following_token] += 1.0
		}

		if _, exist := nm.tokenAppearContext[bigram_tuples]; !exist {
			nm.tokenAppearContext[bigram_tuples] = make([]string, 0)
			nm.tokenAppearContext[bigram_tuples] = append(nm.tokenAppearContext[bigram_tuples], following_token)
		} else {
			nm.tokenAppearContext[bigram_tuples] = append(nm.tokenAppearContext[bigram_tuples], following_token)
		}
	}
}

func (nm *BiGramModel) GetProabilityWeight() BigramProabilityCollections {
	if nm.Pretrained {
		return nm.TokenProabilityWeight
	}

	nm.calucateFullProability()
	return nm.TokenProabilityWeight
}

func (nm *BiGramModel) GetSize() int {
	return len(nm.tokenProability)
}

func (nm *BiGramModel) calculateTokenWeight(prevbitoken BiGramTokenTuple, nexttoken string) float64 {
	context_total := float64(len(nm.tokenAppearContext[prevbitoken]))
	proability_total := nm.tokenProability[prevbitoken][nexttoken]

	return proability_total / context_total
}

func (nm *BiGramModel) calucateFullProability() {
	fmt.Println("Calculating Weight... \nSize: ", len(nm.tokenProability))
	for keyprevtoken, internal_tokenmap := range nm.tokenProability {
		for keynexttoken := range internal_tokenmap {
			nm.TokenProabilityWeight[keyprevtoken] = make(map[string]float64)
			nm.TokenProabilityWeight[keyprevtoken][keynexttoken] = nm.calculateTokenWeight(keyprevtoken, keynexttoken)
		}
	}
}
