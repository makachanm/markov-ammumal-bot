package core

type tokenTuple struct {
	Token      string
	Proability float64
}

type tokenTupleGroups []tokenTuple

func (tkg tokenTupleGroups) Len() int {
	return len(tkg)
}

func (tkg tokenTupleGroups) Swap(i, j int) {
	tkg[i], tkg[j] = tkg[j], tkg[i]
}

func (tkg tokenTupleGroups) Less(i, j int) bool {
	return tkg[i].Proability < tkg[j].Proability
}

func freqInArray(x string, a []string) int {
	var q int = 0
	for _, v := range a {
		if v == x {
			q++
		}
	}

	return q
}

func SerializeBigram(cl BiGramModel) SerializedBigramProabilityCollection {
	seriz := SerializedBigramProabilityCollection{
		Datas: make([]SerializedBigramProabilityCollectionAtom, 0),
	}

	for key, val := range cl.GetProabilityWeight() {
		obj := SerializedBigramProabilityCollectionAtom{
			FirstToken: key.FirstToken,
			NextToken:  key.NextToken,
			Index:      val,
		}

		seriz.Datas = append(seriz.Datas, obj)
	}

	return seriz
}

func UnserializeBigram(cl SerializedBigramProabilityCollection) BiGramModel {
	unsriz := NewBiGramModel()
	unsriz.Pretrained = true

	for _, val := range cl.Datas {
		unsriz.TokenProabilityWeight[BiGramTokenTuple{FirstToken: val.FirstToken, NextToken: val.NextToken}] = val.Index
	}

	return unsriz
}
