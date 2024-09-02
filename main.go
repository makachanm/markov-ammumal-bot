package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"randomsentensbot/core"
	"randomsentensbot/misskey"
	"strings"
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

		if len(config.TwitterData) != 0 {
			core.LoadTwitter(config.TwitterData)
		}

		if len(config.MisskeyData) != 0 {
			core.LoadMisskey(config.MisskeyData)
		}
	} else {
		fmt.Println("Running with Pretrain")

		core.LoadPretrain(config.Pretrain.DataPath)
	}

	predictr = core.GetPredictr()

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

makeText:
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

	if _, e := url.ParseRequestURI(presult.Result); e == nil {
		goto makeText
	}

	text := presult.Result

	text = strings.ReplaceAll(text, "&amp;", "&")
	text = strings.ReplaceAll(text, "&lt;", "<")
	text = strings.ReplaceAll(text, "&gt;", ">")
	text = strings.ReplaceAll(text, "&quot;", "\"")
	text = strings.ReplaceAll(text, "&apos;", "'")

	fmt.Println(presult)
	mk.SendNote(text, vrange)
}

func pretrain(c Config, name string) {
	unifile, fuerr := os.Create(name)
	if fuerr != nil {
		panic(fuerr)
	}
	defer unifile.Close()

	core.PreanalysisData(c.TwitterData, unifile)
}
