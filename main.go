package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"randomsentensbot/core"
	"randomsentensbot/misskey"
	"time"
)

func main() {
	configpath := flag.String("c", "./config.json", "path of configuration file")
	runpretrain := flag.Bool("pt", false, "run pretrain and save to file")
	pretrainname := flag.String("ptf", "data.json", "name of pretrained file")

	flag.Parse()

	config := ReadConfig(*configpath)

	if *runpretrain {
		pretrain(config, *pretrainname)
		return
	}

	var predictr core.PredictionGenerator

	if !config.Pretrain.UsePretrain {
		fmt.Println("Running with Hotload")
		predictr = core.Predictor(config.DataPath)
	} else {
		fmt.Println("Running with Pretrain")

		d, pderr := os.ReadFile(config.Pretrain.DataPath)
		if pderr != nil {
			panic(pderr)
		}

		predictr = core.PreloadPredictor(d)
	}

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

	fmt.Println(presult)
	mk.SendNote(presult.Result, vrange)
}

func pretrain(c Config, name string) {
	unifile, fuerr := os.Create(name)
	if fuerr != nil {
		panic(fuerr)
	}
	defer unifile.Close()

	core.PreanalysisData(c.DataPath, unifile)
}
