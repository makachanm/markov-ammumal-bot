package misskey

type ViewRange string

const PUBLIC ViewRange = "public"
const HOME ViewRange = "home"
const PRIVATE ViewRange = "followers"
const DIRECT ViewRange = "specified"

type MkNote struct {
	Token   string  `json:"i"`
	Content string  `json:"text"`
	CW      *string `json:"cw"`

	ReplyID *string `json:"replyId"`

	LocalOnly  bool      `json:"localOnly"`
	ShareRange ViewRange `json:"visibility"`
}
