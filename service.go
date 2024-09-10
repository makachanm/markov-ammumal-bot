package main

import (
	"context"
	"fmt"
	"randomsentensbot/core"
	"randomsentensbot/misskey"
	"randomsentensbot/service"
	"randomsentensbot/utils"

	"github.com/reugn/go-quartz/quartz"
)

type InteralServices struct {
	Scheduler     quartz.Scheduler
	Configuration utils.Config

	AutoGen       service.AutoGenerationService
	QuestionReply service.QuestionReplierService
}

type InternalServiceArguments struct {
	Config             utils.Config
	Misskey            misskey.Misskey
	PredictonGenerator core.PredictionGenerator

	UniModel core.UniGramModel
}

func NewInternalServiceManager(isa InternalServiceArguments) InteralServices {
	return InteralServices{
		Scheduler:     quartz.NewStdScheduler(),
		Configuration: isa.Config,

		AutoGen:       service.NewAutoGenerationService(isa.PredictonGenerator, isa.Misskey, isa.Config),
		QuestionReply: service.NewQuestionReplierService(isa.UniModel, isa.PredictonGenerator, isa.Misskey),
	}
}

func (is *InteralServices) InitService() {
	ctx := context.Background()

	is.Scheduler.Start(ctx)
	fmt.Println("Scheduler stated")

	//is.Scheduler.ScheduleJob(quartz.NewJobDetail(is.AutoGen, quartz.NewJobKey("AutoGeneration")), quartz.NewSimpleTrigger(time.Second*20))
	is.Scheduler.ScheduleJob(quartz.NewJobDetail(is.QuestionReply, quartz.NewJobKey("QuetionReplier")), quartz.NewRunOnceTrigger(0))

	is.Scheduler.Wait(ctx)
}
