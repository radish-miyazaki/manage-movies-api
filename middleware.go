package main

import (
	"errors"
	"github.com/pascaldekloe/jwt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (app *application) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")  // 全オリジンを許可（ドメイン）
		w.Header().Set("Access-Control-Allow-Headers", "*") // 全Headerを許可（Content-Typeなど）
		w.Header().Set("Access-Control-Allow-Methods", "*") // 全Methodを許可（GETなど）

		// 次のHandlerにチェーンする
		next.ServeHTTP(w, r)
	})
}

func (app *application) checkToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")

		// Headerに含まれるAuthorizationの値を変数に切り出す
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			// could set an anonymous user
			app.errorJSON(w, errors.New("unauhtorization"))
			return
		}

		headerParts := strings.Split(authHeader, " ")

		// 値が正しい構成になっているか
		if len(headerParts) != 2 {
			app.errorJSON(w, errors.New("invalid auth header"))
			return
		}

		// 先頭に"Bearer"が含まれるか
		if headerParts[0] != "Bearer" {
			app.errorJSON(w, errors.New("unauthorized - no bearer"))
			return
		}

		// tokenの値のみを切り出す
		token := headerParts[1]

		// tokenが不正なものでないかチェック
		claims, err := jwt.HMACCheck([]byte(token), []byte(app.config.jwt.secret))
		if err != nil {
			app.errorJSON(w, errors.New("unauthorized - failed hmac check"), http.StatusForbidden)
			return
		}

		// tokenが有効期限内かチェック
		if !claims.Valid(time.Now()) {
			app.errorJSON(w, errors.New("unauthorized - token expired"), http.StatusForbidden)
			return
		}

		// tokenの発行者が正しいかチェック
		if claims.Issuer != "my_domain.com" {
			app.errorJSON(w, errors.New("unauthorized - invalid issuer"), http.StatusForbidden)
			return
		}

		// tokenの想定利用者が正しいかチェック
		if !claims.AcceptAudience("my_domain.com") {
			app.errorJSON(w, errors.New("unauthorized - invalid audience"), http.StatusForbidden)
			return
		}

		userID, err := strconv.ParseInt(claims.Subject, 10, 64)
		if err != nil {
			app.errorJSON(w, errors.New("unauthorized"), http.StatusForbidden)
			return
		}

		log.Println("UserID is " + strconv.Itoa(int(userID)))

		// 次のHandlerにチェーンする
		next.ServeHTTP(w, r)
	})
}
