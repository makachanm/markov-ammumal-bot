package utils

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
	Cron     cronSettings   `json:"cron"`

	UseReply bool `json:"usereply"`
}

type configPreTrain struct {
	UsePretrain bool   `json:"usepretrain"`
	DataPath    string `json:"path"`
}

type cronSettings struct {
	UseCron bool   `json:"usecron"`
	Crontab string `json:"crontab"`
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
