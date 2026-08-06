//go:build js && wasm

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"syscall/js"

	appContext "github.com/DenysSkobalo/denysskobalodev.space/internal/context"
)

var muxHandler http.Handler

func main() {
	muxHandler = setupRouter()
	js.Global().Set("handleHttpRequest", js.FuncOf(handleHttpRequest))
	select {}
}

func handleHttpRequest(this js.Value, args []js.Value) (resp interface{}) {
	defer func() {
		if r := recover(); r != nil {
			errResp := map[string]interface{}{
				"status": 500,
				"body":   fmt.Sprintf(`{"error":"Internal Wasm Panic","details":%q}`, fmt.Sprintf("%v", r)),
			}
			respBytes, _ := json.Marshal(errResp)
			resp = string(respBytes)
		}
	}()

	if len(args) == 0 {
		return `{"status":400,"body":"{\"error\":\"Missing request payload\"}"}`
	}

	jsReq := args[0]
	method := jsReq.Get("method").String()
	urlStr := jsReq.Get("url").String()
	bodyStr := jsReq.Get("body").String()

	adminHash := ""
	salt := "denysskobalo_unique_salt"

	if jsEnv := jsReq.Get("env"); !jsEnv.IsUndefined() && !jsEnv.IsNull() {
		if h := jsEnv.Get("ADMIN_HASH"); !h.IsUndefined() {
			adminHash = strings.TrimSpace(h.String())
		}
		if s := jsEnv.Get("SALT"); !s.IsUndefined() {
			salt = strings.TrimSpace(s.String())
		}
	}

	req := httptest.NewRequest(method, urlStr, strings.NewReader(bodyStr))

	ctx := context.WithValue(req.Context(), appContext.AdminHashKey, adminHash)
	ctx = context.WithValue(ctx, appContext.SaltKey, salt)
	req = req.WithContext(ctx)

	rec := httptest.NewRecorder()
	muxHandler.ServeHTTP(rec, req)

	respMap := map[string]interface{}{
		"status": rec.Code,
		"body":   rec.Body.String(),
	}

	respBytes, err := json.Marshal(respMap)
	if err != nil {
		return `{"status":500,"body":"{\"error\":\"Failed to serialize WASM response\"}"}`
	}

	return string(respBytes)
}
