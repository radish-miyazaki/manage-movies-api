package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/pascaldekloe/jwt"
	"github.com/radish-miyazaki/manage-movies-api/models"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

// モック用ユーザ
var validUser = models.User{
	ID:       10,
	Email:    "me@example.com",
	Password: "$2a$12$FEwS8U1cK/hzG6oNOBGp5OL0JUKQSP6FrMsBjv8wfbJrWNdaGUwOC",
}

// Credentials 認証情報
type Credentials struct {
	UserName string `json:"email"`
	Password string `json:"password"`
}

// Login ログイン認証を行う
func (app *application) Login(w http.ResponseWriter, r *http.Request) {
	var cs Credentials

	// リクエストをデコードし、構造体に格納
	err := json.NewDecoder(r.Body).Decode(&cs)
	if err != nil {
		app.errorJSON(w, errors.New("unauthorized"))
		return
	}

	hashedPassword := validUser.Password

	// メールアドレスチェック
	if cs.UserName != validUser.Email {
		app.errorJSON(w, errors.New("unauthorized"))
		return
	}

	// リクエストのパスワードとDBで保管しているハッシュ化したパスワードを比較
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(cs.Password))
	if err != nil {
		app.errorJSON(w, errors.New("unauthorized"))
		return
	}

	// jwtのユーザー情報であるclaimsを作成
	var claims jwt.Claims
	claims.Subject = fmt.Sprint(validUser.ID)
	claims.Issued = jwt.NewNumericTime(time.Now())
	claims.NotBefore = jwt.NewNumericTime(time.Now())
	claims.Expires = jwt.NewNumericTime(time.Now().Add(24 * time.Hour))
	claims.Issuer = "my_domain.com"
	claims.Audiences = []string{"my_domain.com"}

	// 作成したclaimを暗号化し、tokenを作成
	jwtBytes, err := claims.HMACSign(jwt.HS256, []byte(app.config.jwt.secret))
	if err != nil {
		app.errorJSON(w, errors.New("error login"))
		return
	}

	// 暗号化したtokenをレスポンスに付与
	if err = app.writeJSON(w, http.StatusOK, string(jwtBytes), "response"); err != nil {
		app.errorJSON(w, errors.New("error login"))
		return
	}
}
