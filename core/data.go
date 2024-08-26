package core

import (
	"encoding/json"
	"fmt"
	"os"
)

type DNotes struct {
	Text string `json:"text"`
}

func Predictor(path string) PredictionGenerator {
	bd, fe := os.ReadFile(path)
	if fe != nil {
		panic(fe)
	}

	var nd []DNotes
	xe := json.Unmarshal(bd, &nd)
	if xe != nil {
		panic(xe)
	}

	sm := NewUniGramModel()
	bm := NewBiGramModel()

	fmt.Println("Updating Data....")

	for _, av := range nd {
		bm.Update(av.Text)
	}
	for _, av := range nd {
		sm.Update(av.Text)
	}

	tg := NewPredictionGenerator(sm, bm)
	return tg
}
