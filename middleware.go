package main

import "net/http"

func (app *application) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")  // 全オリジンを許可（ドメイン）
		w.Header().Set("Access-Control-Allow-Headers", "*") // 全Headerを許可（Content-Typeなど）
		w.Header().Set("Access-Control-Allow-Methods", "*") // 全Methodを許可（GETなど）

		// 次のHandlerにチェーンする
		next.ServeHTTP(w, r)
	})
}
