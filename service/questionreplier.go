package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"randomsentensbot/core"
	"randomsentensbot/misskey"
)

type QuestionReplierService struct {
	desc string

	uniModel  core.UniGramModel
	generator core.PredictionGenerator
	mk        misskey.Misskey
}

func NewQuestionReplierService(uc core.UniGramModel, gen core.PredictionGenerator, mks misskey.Misskey) QuestionReplierService {
	return QuestionReplierService{
		uniModel:  uc,
		generator: gen,
		mk:        mks,
	}
}

func (qrs QuestionReplierService) Description() string {
	return qrs.desc
}

func (qrs QuestionReplierService) Execute(ctx context.Context) error {
	var mux http.ServeMux = *http.NewServeMux()
	mux.HandleFunc("POST /hook", qrs.handleHook)

	http.ListenAndServe(":3000", &mux)

	return nil
}

func (qrs QuestionReplierService) handleHook(res http.ResponseWriter, req *http.Request) {
	var hookData WebhookData = *new(WebhookData)

	rawNoteData, readerror := io.ReadAll(req.Body)
	if readerror != nil {
		log.Fatal("readerror:", readerror.Error())
	}
	json.Unmarshal(rawNoteData, &hookData)

	if !(hookData.HookType == TYPE_MENTION || hookData.HookType == TYPE_REPLY) {
		fmt.Println("request: skip")
		res.WriteHeader(http.StatusNotAcceptable)
		res.Write([]byte("failed"))

		return
	} else {
		res.Write([]byte("accepted"))
	}

	extractor := core.NewImportantExtractor(qrs.uniModel)

	extracted := extractor.Extract(hookData.Body.Note.Text)

	var generated core.PredictionResult
	for _, token := range extracted {
		result := qrs.generator.PredictSeq(token.Token, 0)
		if len(result.Seq) > 1 {
			generated = result
			break
		}
	}

	qrs.mk.SendReply(hookData.Body.Note.NoteID, generated.Result, hookData.Body.Note.Visibility)
}
