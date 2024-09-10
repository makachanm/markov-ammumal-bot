package service

import "randomsentensbot/misskey"

type HookType string

const TYPE_MENTION HookType = "mention"
const TYPE_REPLY HookType = "reply"
const TYPE_FOLLOWED HookType = "followed"
const TYPE_FOLLOW HookType = "follow"

type WebhookData struct {
	ServerURL string   `json:"server"`
	HookType  HookType `json:"type"`

	Body hookBody `json:"body,omitempty"`
}

type hookBody struct {
	Note MisskeyHookNote `json:"note,omitempty"`
}

type MisskeyHookNote struct {
	NoteID string `json:"id"`
	UserID string `json:"userId"`

	LocalOnly  bool              `json:"localOnly"`
	Visibility misskey.ViewRange `json:"visibility"`
	Text       string            `json:"text"`
	CW         string            `json:"cw"`
}
