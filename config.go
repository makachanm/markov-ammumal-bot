package main

import (
	"encoding/json"
	"os"
)

type Config struct {
	MisskeyToken  string   `json:"mktoken"`
	MisskeyServer string   `json:"mkserver"`
	ViewRange     string   `json:"range"`
	StartTopic    []string `json:"starttopic"`
	MisskeyData   []string `json:"misskeyData"`
	TwitterData   []string `json:"twitterData"`

	Pretrain configPreTrain `json:"pretrain"`
}

type configPreTrain struct {
	UsePretrain bool   `json:"usepretrain"`
	DataPath    string `json:"path"`
}

func ReadConfig(path string) Config {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	cf := *new(Config)
	jerr := json.Unmarshal(data, &cf)
	if jerr != nil {
		panic(jerr)
	}

	return cf
}
