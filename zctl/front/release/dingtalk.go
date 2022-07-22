package release

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"net/url"
)

const URL = "https://oapi.dingtalk.com/robot/send?access_token="
const secret = "SECbd48f4e2fc5a85dceb453cfa917db016ab3bb21f348cde132a1f26bd4149e56b"

type DingTalk struct {
	token  string
	secret string
}

func DingTalkNew(secret, token string) *DingTalk {
	return &DingTalk{
		secret: secret,
		token:  token,
	}
}

func (d *DingTalk) Talk(title, content string, AtMobiles []string, AtUserIds []string, IsAtAll bool) error {
	t := new(dMarkdown)
	t.Text = struct {
		Title   string `json:"title"`
		Content string `json:"text"`
	}{Title: title, Content: content}
	t.At = struct {
		AtMobiles []string `json:"atMobiles"`
		AtUserIds []string `json:"atUserIds"`
		IsAtAll   bool     `json:"isAtAll"`
	}(struct {
		AtMobiles []string
		AtUserIds []string
		IsAtAll   bool
	}{AtMobiles: AtMobiles, AtUserIds: AtUserIds, IsAtAll: IsAtAll})
	t.Msgtype = "markdown"
	if jsonByte, err := json.Marshal(t); err != nil {
		return err
	} else if _, err := d.send(jsonByte, URL+d.token); err != nil {
		return err
	}
	return nil
}

// Text 文本json
type dMarkdown struct {
	At struct {
		AtMobiles []string `json:"atMobiles"`
		AtUserIds []string `json:"atUserIds"`
		IsAtAll   bool     `json:"isAtAll"`
	} `json:"at"`
	Text struct {
		Title   string `json:"title"`
		Content string `json:"text"`
	} `json:"markdown"`
	Msgtype string `json:"msgtype"`
}

//Link Link型json
type dLink struct {
	Msgtype string `json:"msgtype"`
	Link    struct {
		Text       string `json:"text"`
		Title      string `json:"title"`
		PicUrl     string `json:"picUrl"`
		MessageUrl string `json:"messageUrl"`
	} `json:"link"`
}

//MD Markdown型json
type dMD struct {
	Msgtype  string `json:"msgtype"`
	Markdown struct {
		Title string `json:"title"`
		Text  string `json:"text"`
	} `json:"markdown"`
	At struct {
		AtMobiles []string `json:"atMobiles"`
		AtUserIds []string `json:"atUserIds"`
		IsAtAll   bool     `json:"isAtAll"`
	} `json:"at"`
}

// AActionCard 整体跳转ActionCard类型
type dAActionCard struct {
	ActionCard struct {
		Title          string `json:"title"`
		Text           string `json:"text"`
		BtnOrientation string `json:"btnOrientation"`
		SingleTitle    string `json:"singleTitle"`
		SingleURL      string `json:"singleURL"`
	} `json:"actionCard"`
	Msgtype string `json:"msgtype"`
}

// DActionCard 独立跳转ActionCard类型
type dDActionCard struct {
	Msgtype    string `json:"msgtype"`
	ActionCard struct {
		Title          string `json:"title"`
		Text           string `json:"text"`
		BtnOrientation string `json:"btnOrientation"`
		Btns           []struct {
			Title     string `json:"title"`
			ActionURL string `json:"actionURL"`
		} `json:"btns"`
	} `json:"actionCard"`
}

// ErrorReport 返回的错误
type dErrorReport struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

const httpTimoutSecond = time.Duration(30) * time.Second

func (d *DingTalk) send(message []byte, pushURL string) (*dErrorReport, error) {
	timestamp := time.Now().UnixMilli()
	pushURL = pushURL + "&timestamp=" + fmt.Sprintf("%d", timestamp)
	sign := hmacSha256(fmt.Sprintf("%d", timestamp)+"\n"+d.secret, d.secret)
	pushURL = pushURL + "&sign=" + sign
	res := &dErrorReport{}

	reqBytes := message

	req, err := http.NewRequest(http.MethodPost, pushURL, bytes.NewReader(reqBytes))
	if err != nil {
		return res, err
	}
	req.Header.Add("Accept-Charset", "utf8")
	req.Header.Add("Content-Type", "application/json")

	client := new(http.Client)
	client.Timeout = httpTimoutSecond
	resp, err := client.Do(req)
	if err != nil {
		return res, err
	}
	defer resp.Body.Close()

	resultByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return res, err
	}

	err = json.Unmarshal(resultByte, &res)
	if err != nil {
		return res, fmt.Errorf("unmarshal http response body from json error = %w", err)
	}

	if res.Errcode != 0 {
		return res, fmt.Errorf("send message to dingtalk error = %s", res.Errmsg)
	}

	return res, nil
}
func hmacSha256(data string, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	return url.PathEscape(base64.StdEncoding.EncodeToString(h.Sum(nil)))
}
