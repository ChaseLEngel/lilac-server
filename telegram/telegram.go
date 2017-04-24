package telegram

import (
	"fmt"
	"net/http"
	"net/url"
)

type Telegram struct {
	token        string
	tokenizedUrl string
}

func New(token string) Telegram {
	telegram := Telegram{}
	telegram.token = token
	telegram.tokenizedUrl = fmt.Sprintf("https://api.telegram.org/bot%v/sendMessage", token)
	return telegram
}

func (t Telegram) SendMessage(chatId, message string) error {
	values := url.Values{"chat_id": {chatId}, "text": {message}}
	resp, err := http.PostForm(t.tokenizedUrl, values)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("%v", resp.Status)
	}
	return nil
}
