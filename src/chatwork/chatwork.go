package chatwork

import (
	"gobot/src/env"
	"gobot/src/gae"
	"net/http"
	"net/url"
)

// SendMessage sends message to chatwork
func SendMessage(r *http.Request, roomID, body string) error {
	// Create request
	req, err := http.NewRequest("POST", "https://api.chatwork.com/v2/rooms/"+roomID+"/messages", nil)
	if err != nil {
		return err
	}
	req.Header.Add("X-ChatWorkToken", env.ChatWorkToken)

	// クエリを組み立て
	values := url.Values{}   // url.Valuesオブジェクト生成
	values.Add("body", body) // key-valueを追加
	req.URL.RawQuery = values.Encode()

	resp, err := gae.Do(r, req)
	if err != nil {
		return err
	}

	// 関数を抜ける際に必ずresponseをcloseするようにdeferでcloseを呼ぶ
	defer resp.Body.Close()

	return nil
}
