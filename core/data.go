package core

import (
	"encoding/json"
	"log"
	"os"

	"github.com/tidwall/gjson"
)

type modelData struct {
	BiModel  SerializedBigramProabilityCollection `json:"bi"`
	UniModel UnigramProabilityCollections         `json:"uni"`
}

var (
	sm UniGramModel
	bm BiGramModel
)

func init() {
	sm = NewUniGramModel()
	bm = NewBiGramModel()
}

func generalLoader(paths []string, textPath string) {
	var texts []string = make([]string, 0)

	for _, data := range paths {
		bytes, err := os.ReadFile(data)
		if err != nil {
			panic(err)
		}

		values := gjson.Get(string(bytes), textPath)
		for _, value := range values.Array() {
			texts = append(texts, value.String())
		}
	}

	for _, text := range texts {
		bm.Update(text)
		sm.Update(text)
	}
}

func LoadMisskey(paths []string) {
	generalLoader(paths, `notes.#(visibility!="specified").text`)
}

func LoadTwitter(paths []string) {
	generalLoader(paths, `#.tweet.full_text`)
}

func LoadPretrain(path string) {
	ubmx := modelData{}

	file, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(file, &ubmx)
	if err != nil {
		panic(err)
	}

	sm.TokenProabilityWeight = ubmx.UniModel
	sm.Pretrained = true
	bm = UnserializeBigram(ubmx.BiModel)
}

func GetPredictr() PredictionGenerator {
	return NewPredictionGenerator(sm, bm)
}

func PreanalysisData(paths []string, writer *os.File) {
	generalLoader(paths, `notes.#(visibility!="specified").text`)

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
