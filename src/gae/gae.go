package gae

import (
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
)

// Client get GAE http.Client
func Do(r, req *http.Request) (*http.Response, error) {
	// Doメソッドでリクエストを投げる
	// http.Response型のポインタ（とerror）が返ってくる
	ctx := appengine.NewContext(r)
	client := urlfetch.Client(ctx)

	return client.Do(req)
}
