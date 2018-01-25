package app

import (
	"encoding/json"
	"fmt"
	"gobot/src/chatwork"
	"gobot/src/env"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func init() {

	env.Load()

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/time", timeHandler)
	http.HandleFunc("/send", sendHandler)
	http.HandleFunc("/mention", mentionHandler)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	// rootかどうかチェック
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// response
	fmt.Fprint(w, "Instance is alive!")
}

func sendHandler(w http.ResponseWriter, r *http.Request) {
	//Validate request
	if r.Method != "GET" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// query params
	body := r.URL.Query().Get("message")
	if len(body) < 1 {
		fmt.Fprintln(w, "Url Param 'message' is missing")
		return
	}

	// 送信
	if err := chatwork.SendMessage(r, env.RoomID, body); err != nil {
		fmt.Fprintln(w, err)
	}

	// response
	fmt.Fprintln(w, "Send message!")
}

func timeHandler(w http.ResponseWriter, r *http.Request) {
	//Validate request
	if r.Method != "GET" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// メッセージ作成
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	body := strconv.Itoa(time.Now().In(jst).Hour()) + "時ですよ"

	// 送信
	if err := chatwork.SendMessage(r, env.RoomID, body); err != nil {
		fmt.Fprintln(w, err)
	}

	// response
	fmt.Fprintln(w, "Hello, Chatwork!")
}

/** JSONデコード用に構造体定義 */
type cwWebhook struct {
	WebhookSettingID string `json:"webhook_setting_id"`
	WebhookEventType string `json:"webhook_event_type"`
	WebhookEventTime int    `json:"webhook_event_time"`
	WebhookEvent     struct {
		FromAccountID int    `json:"from_account_id"`
		ToAccountID   int    `json:"to_account_id"`
		RoomID        int    `json:"room_id"`
		MessageID     string `json:"message_id"`
		Body          string `json:"body"`
		SendTime      int    `json:"send_time"`
		UpdateTime    int    `json:"update_time"`
	} `json:"webhook_event"`
}

func mentionHandler(w http.ResponseWriter, r *http.Request) {
	// Validate request
	if r.Method != "POST" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// To allocate slice for request body
	length, err := strconv.Atoi(r.Header.Get("Content-Length"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Read body data to parse json
	body := make([]byte, length)
	length, err = r.Body.Read(body)
	if err != nil && err != io.EOF {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// parse json
	var jsonBody cwWebhook
	err = json.Unmarshal(body[:length], &jsonBody)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Create message
	var arr = []string{}
	roomID := strconv.Itoa(jsonBody.WebhookEvent.RoomID)
	arr = append(arr, "[rp aid="+strconv.Itoa(jsonBody.WebhookEvent.FromAccountID)+" to="+roomID+"-"+jsonBody.WebhookEvent.MessageID+"]")
	for _, v := range regexp.MustCompile("\r\n|\n\r|\n|\r").Split(jsonBody.WebhookEvent.Body, -1) {
		if !strings.Contains(v, "[To:"+strconv.Itoa(jsonBody.WebhookEvent.ToAccountID)+"]") {
			arr = append(arr, v)
		}
	}
	message := strings.Join(arr, "\n")

	// Send message
	err = chatwork.SendMessage(r, env.RoomID, message)
	if err != nil {
		fmt.Fprintln(w, err)
	}

	// Response
	fmt.Fprintln(w, "Hello, chatwork!")
}
