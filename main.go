package main

import (
	"randomsentensbot/core"
	"randomsentensbot/misskey"
)

func main() {
	config := ReadConfig()
	predictr := core.Predictor(config.DataPath)

	var vrange misskey.ViewRange

	switch config.ViewRange {
	case "public":
		vrange = misskey.PUBLIC
	case "home":
		vrange = misskey.HOME
	case "private":
		vrange = misskey.PRIVATE
	default:
		vrange = misskey.HOME
	}

	mk := misskey.NewMisskeyTools(config.MisskeyToken, config.MisskeyServer)
	presult := predictr.PredictSeq(config.StartTopic, 0)

	mk.SendNote(presult.Result, vrange)
}
