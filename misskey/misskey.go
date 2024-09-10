package misskey

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type Misskey struct {
	token     string
	serverurl string
}

func NewMisskeyTools(token string, serverurl string) Misskey {
	return Misskey{
		token:     token,
		serverurl: serverurl,
	}
}

func (m *Misskey) SendNote(content string, vrange ViewRange) {
	ctx := MkNote{
		Content:    content,
		ShareRange: vrange,
		LocalOnly:  false,

		Token: m.token,
	}

	x, _ := json.Marshal(ctx)
	body := bytes.NewReader(x)

	url, _ := url.Parse(m.serverurl)
	xurl := url.JoinPath("api", "notes", "create")

	fmt.Println(xurl.String())

	reqs, xerr := http.NewRequest(http.MethodPost, xurl.String(), body)
	reqs.Header.Set("Content-Type", "application/json")

	if xerr != nil {
		fmt.Println(xerr)
		return
	}

	client := http.Client{}
	_, rerr := client.Do(reqs)
	if rerr != nil {
		fmt.Println(rerr)
		return
	}
}

func (m *Misskey) SendReply(replyid string, content string, vrange ViewRange) {
	ctx := MkNote{
		ReplyID: &replyid,

		Content:    content,
		ShareRange: vrange,

		Token: m.token,
	}

	x, _ := json.Marshal(ctx)
	body := bytes.NewReader(x)

	url, _ := url.Parse(m.serverurl)
	xurl := url.JoinPath("api", "notes", "create")

	fmt.Println(xurl.String())

	reqs, xerr := http.NewRequest(http.MethodPost, xurl.String(), body)
	reqs.Header.Set("Content-Type", "application/json")

	if xerr != nil {
		fmt.Println(xerr)
		return
	}

	client := http.Client{}
	_, rerr := client.Do(reqs)
	if rerr != nil {
		fmt.Println(rerr)
		return
	}
}
