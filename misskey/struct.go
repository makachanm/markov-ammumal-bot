package misskey

type ViewRange string

const PUBLIC ViewRange = "public"
const HOME ViewRange = "home"
const PRIVATE ViewRange = "followers"

type MkNote struct {
	Token   string `json:"i"`
	Content string `json:"text"`

	LocalOnly  bool      `json:"localOnly"`
	ShareRange ViewRange `json:"visibility"`
}
