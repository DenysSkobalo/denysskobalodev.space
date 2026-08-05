//go:build js && wasm

package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"syscall/js"
)

var muxHandler http.Handler

func main() {
	muxHandler = setupRouter()

	// Експортуємо JS-функцію у глобальний контекст V8 Isolate
	js.Global().Set("handleHttpRequest", js.FuncOf(handleHttpRequest))

	// Блокуємо головну горутину Go WASM
	select {}
}

func handleHttpRequest(this js.Value, args []js.Value) any {
	if len(args) == 0 {
		return `{"status":400,"body":"{\"error\":\"Missing request argument\"}"}`
	}

	jsReq := args[0]
	method := jsReq.Get("method").String()
	urlStr := jsReq.Get("url").String()
	bodyStr := jsReq.Get("body").String()

	req := httptest.NewRequest(method, urlStr, strings.NewReader(bodyStr))
	rec := httptest.NewRecorder()

	muxHandler.ServeHTTP(rec, req)

	respMap := map[string]any{
		"status": rec.Code,
		"body":   rec.Body.String(),
	}

	respBytes, err := json.Marshal(respMap)
	if err != nil {
		return `{"status":500,"body":"{\"error\":\"Failed to marshal WASM response\"}"}`
	}

	return string(respBytes)
}
