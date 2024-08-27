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
	DataPath      string   `json:"dataname"`
}

func ReadConfig() Config {
	data, err := os.ReadFile("./config.json")
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
