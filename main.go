package main

import (
	"math/rand"
	"randomsentensbot/core"
	"randomsentensbot/misskey"
	"time"
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

	rand.Seed(time.Now().Unix())
	topic := config.StartTopic[rand.Intn(len(config.StartTopic))]

	if topic == "random" {
		pick := func(length int, dict core.UnigramProabilityCollections) string {
			rndn := rand.Intn(length)
			for key := range dict {
				if rndn == 0 {
					return key
				}
				rndn--
			}
			panic("unreachable!")
		}

		topic = pick(len(predictr.UniModelProb), predictr.UniModelProb)
	}

	presult := predictr.PredictSeq(topic, 0)

	mk.SendNote(presult.Result, vrange)
}
