package main

import (
	"flag"
	"fmt"
	"os"
	"randomsentensbot/core"
	"randomsentensbot/misskey"
	"randomsentensbot/service"
	"randomsentensbot/utils"
)

func main() {
	configpath := flag.String("c", "./config.json", "path of configuration file")
	runpretrain := flag.Bool("pt", false, "run pretrain and save to file")
	pretrainname := flag.String("ptf", "data.json", "name of pretrained file")

	headless := flag.Bool("hl", false, "run bot without server things(headless)")

	flag.Parse()

	fmt.Println("Stating...")
	config := utils.ReadConfig(*configpath)

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

	fmt.Println("Core predictor load complete")
	predictr = core.GetPredictr()

	var misskey misskey.Misskey = misskey.NewMisskeyTools(config.MisskeyToken, config.MisskeyServer)

	if !*headless {
		fmt.Println("Starting service...")
		var BotService InteralServices = NewInternalServiceManager(InternalServiceArguments{
			Config:             config,
			PredictonGenerator: predictr,
			Misskey:            misskey,

			UniModel: predictr.UniModel,
		})

		BotService.InitService()
	} else {
		fmt.Println("Running with headless mode")
		service.Generator(config, misskey, predictr)
	}
}

func pretrain(c utils.Config, name string) {
	unifile, fuerr := os.Create(name)
	if fuerr != nil {
		panic(fuerr)
	}
	defer unifile.Close()

	if len(c.TwitterData) != 0 {
		core.LoadTwitter(c.TwitterData)
	}

	if len(c.MisskeyData) != 0 {
		core.LoadMisskey(c.MisskeyData)
	}

	core.PreanalysisData(unifile)
}
