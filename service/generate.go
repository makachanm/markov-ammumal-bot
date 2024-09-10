package service

import (
	"context"
	"fmt"
	"net/url"
	"randomsentensbot/core"
	"randomsentensbot/misskey"
	"randomsentensbot/utils"
	"strings"

	"golang.org/x/exp/rand"
)

type AutoGenerationService struct {
	desc string

	mk        misskey.Misskey
	predictor core.PredictionGenerator
	config    utils.Config
}

func NewAutoGenerationService(pd core.PredictionGenerator, msky misskey.Misskey, conf utils.Config) AutoGenerationService {
	return AutoGenerationService{
		mk:        msky,
		predictor: pd,
		config:    conf,
	}
}

func (ags AutoGenerationService) Description() string {
	return ags.desc
}

func (ags AutoGenerationService) Execute(c context.Context) error {
	Generator(ags.config, ags.mk, ags.predictor)
	return nil
}

func Generator(config utils.Config, mk misskey.Misskey, predictor core.PredictionGenerator) {
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

	topic := config.StartTopic[rand.Intn(len(config.StartTopic))]

	if topic == "random" {
		for {
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

			topic = pick(len(predictor.UniModelProb), predictor.UniModelProb)

			if _, e := url.ParseRequestURI(topic); e != nil {
				break
			}
		}
	}

	var presult core.PredictionResult

	presult = predictor.PredictSeq(topic, 0)

	text := presult.Result

	text = strings.ReplaceAll(text, "&amp;", "&")
	text = strings.ReplaceAll(text, "&lt;", "<")
	text = strings.ReplaceAll(text, "&gt;", ">")
	text = strings.ReplaceAll(text, "&quot;", "\"")
	text = strings.ReplaceAll(text, "&apos;", "'")

	fmt.Println(presult)
	mk.SendNote(text, vrange)
}
