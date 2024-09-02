package core

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type DNotes struct {
	Text       string `json:"text"`
	Visibility string `json:"visibility"`
}

type modelData struct {
	BiModel  SerializedBigramProabilityCollection `json:"bi"`
	UniModel UnigramProabilityCollections         `json:"uni"`
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

	notes = filterNotes(notes)

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


func filterNotes(notes []DNotes) []DNotes {
	var nnotes []DNotes = make([]DNotes, 0)

	for _, av := range notes {
		if av.Visibility != "specified" {
			nnotes = append(nnotes, av)
		}
	}

	return nnotes
}

func PreloadPredictor(file []byte) PredictionGenerator {
	sm := NewUniGramModel()

	ubmx := modelData{}

	jerr := json.Unmarshal(file, &ubmx)
	if jerr != nil {
		panic(jerr)
	}

	sm.TokenProabilityWeight = ubmx.UniModel
	sm.Pretrained = true

	bm := UnserializeBigram(ubmx.BiModel)

	tg := NewPredictionGenerator(sm, bm)
	return tg
}

func PreanalysisData(paths []string, writer *os.File) {
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
  
  notes = filterNotes(notes)

	fmt.Println("Making pre-analysised data....")

	sm := NewUniGramModel()
	bm := NewBiGramModel()

	fmt.Println("Updating Data....")

	for _, av := range notes {
		bm.Update(av.Text)
	}
	for _, av := range notes {
		sm.Update(av.Text)
	}

	model := modelData{
		UniModel: sm.GetProabilityWeight(),
		BiModel:  SerializeBigram(bm),
	}

	data, err := json.Marshal(model)
	if err != nil {
		log.Fatal(err)
	}

	_, ferr := writer.Write(data)
	if ferr != nil {
		panic(err)
	}
}
