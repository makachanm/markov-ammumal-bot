package core

import (
	"encoding/json"
	"fmt"
	"os"
)

type DNotes struct {
	Text string `json:"text"`
}

func Predictor(paths []string) PredictionGenerator {
	var notes []DNotes = make([]DNotes, 0)

	for _, data := range paths {
		bd, fe := os.ReadFile(data)
		if fe != nil {
			panic(fe)
		}

		var nd []DNotes
		xe := json.Unmarshal(bd, &nd)
		if xe != nil {
			panic(xe)
		}

		notes = append(notes, nd...)
	}

	sm := NewUniGramModel()
	bm := NewBiGramModel()

	fmt.Println("Updating Data....")

	for _, av := range notes {
		bm.Update(av.Text)
	}
	for _, av := range notes {
		sm.Update(av.Text)
	}

	tg := NewPredictionGenerator(sm, bm)
	return tg
}
