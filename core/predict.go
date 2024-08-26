package core

import (
	"sort"
	"strings"
)

const TokenThresholdLimit int = 500

func SortUniTokenByProability(input map[string]float64, appears []string, rel UnigramProabilityCollections) tokenTupleGroups {
	sortcandidates := make(tokenTupleGroups, len(input))

	for k, v := range input {
		var relational float64 = 0.0
		var relationaltokens int = 0

		for _, prevtoken := range appears {
			if _, exist := rel[prevtoken][k]; exist {
				relational += rel[prevtoken][k]
				relationaltokens++
			}
		}

		sortcandidates = append(sortcandidates, tokenTuple{
			Token:      k,
			Proability: v + ((relational / float64(relationaltokens)) - (float64(freqInArray(k, appears)))),
		})
	}

	sort.Sort(sort.Reverse(sortcandidates))

	return sortcandidates
}

func SortBiTokenByProability(input map[string]float64, appears []string, rel BigramProabilityCollections) tokenTupleGroups {
	sortcandidates := make(tokenTupleGroups, len(input))

	for k, v := range input {
		sortcandidates = append(sortcandidates, tokenTuple{
			Token:      k,
			Proability: v - (float64(freqInArray(k, appears) / len(appears))),
		})
	}

	sort.Sort(sort.Reverse(sortcandidates))

	return sortcandidates
}

type PredictionGenerator struct {
	BiModel     BiGramModel
	BiModelProb BigramProabilityCollections

	UniModel     UniGramModel
	UniModelProb UnigramProabilityCollections
}

func NewPredictionGenerator(unictx UniGramModel, bictx BiGramModel) PredictionGenerator {
	return PredictionGenerator{
		UniModel:     unictx,
		UniModelProb: unictx.GetProabilityWeight(),

		BiModel:     bictx,
		BiModelProb: bictx.GetProabilityWeight(),
	}
}

func (pg *PredictionGenerator) predictUninModelNext(seq string, currentseq []string) (string, float64) {
	candidates := pg.UniModelProb[seq]
	sorted := SortBiTokenByProability(candidates, currentseq, pg.BiModelProb)

	if len(sorted) <= 0 {
		return ENDTOKEN, 1.0
	} else {
		return sorted[0].Token, sorted[0].Proability
	}
}

func (pg *PredictionGenerator) predictBiModelNext(firstseq string, secondseq string, currentseq []string) (string, float64) {
	candidates := pg.BiModelProb[biTuple(strings.ToLower(firstseq), strings.ToLower(secondseq))]
	sorted := SortBiTokenByProability(candidates, currentseq, pg.BiModelProb)

	if len(sorted) <= 0 {
		return ENDTOKEN, 1.0
	} else {
		return sorted[0].Token, sorted[0].Proability
	}
}

func (pg *PredictionGenerator) predictHybridNext(seq string, currentseq []string) (string, float64) {
	uniptoken, unix := pg.predictUninModelNext(seq, currentseq)
	if _, exist := pg.BiModelProb[biTuple(seq, uniptoken)]; exist {
		return pg.predictBiModelNext(seq, uniptoken, currentseq)
	} else {
		return uniptoken, unix
	}
}

func (pg *PredictionGenerator) PredictSeq(seq string, seqlen int) PredictionResult {
	var seq_length int = 0
	var seq_nolimit bool = false

	if seqlen == 0 {
		seq_nolimit = true
	} else {
		seq_length = seqlen
	}

	//total predicted strings
	var predicts = make([]string, 0)
	predicts = append(predicts, seq)

	//selection of tokens
	var predict_seq_selection = make([]predictionResultSeqAtom, 0)

	//prev token
	beforePredict := seq

	if !seq_nolimit {
		for i := 0; i < seq_length; i++ {
			predictseq, proability := pg.predictHybridNext(beforePredict, predicts)

			if predictseq == ENDTOKEN {
				break
			}

			beforePredict = predictseq
			predicts = append(predicts, predictseq)

			predict_seq_selection = append(predict_seq_selection, predictionResultSeqAtom{
				Token:      predictseq,
				Proability: proability,
			})
		}
	} else {
		for {
			predictseq, proability := pg.predictHybridNext(beforePredict, predicts)

			if predictseq == ENDTOKEN {
				break
			}

			beforePredict = predictseq
			predicts = append(predicts, predictseq)

			predict_seq_selection = append(predict_seq_selection, predictionResultSeqAtom{
				Token:      predictseq,
				Proability: proability,
			})

			if len(predicts) > TokenThresholdLimit {
				break
			}
		}
	}

	return PredictionResult{
		Result: strings.Join(predicts, " "),
		Seq:    predict_seq_selection,
	}
}
