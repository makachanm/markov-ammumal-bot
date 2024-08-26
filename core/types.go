package core

type UnigramProabilityCollections map[string]map[string]float64
type uniGramContextTokenAppears map[string][]string

type BigramProabilityCollections map[BiGramTokenTuple]map[string]float64
type biGramContextTokenAppears map[BiGramTokenTuple][]string

type BiGramTokenTuple struct {
	FirstToken string
	NextToken  string
}

type PredictionResult struct {
	Result string
	Seq    []predictionResultSeqAtom
}

type predictionResultSeqAtom struct {
	Token      string
	Proability float64
}
