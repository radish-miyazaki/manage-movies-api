package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/radish-miyazaki/manage-movies-api/models"
	"log"
	"net/http"
	"os"
	"time"
)

// version ... application version
const version = "1.0"

// config ... application configuration
type config struct {
	port int
	env  string
	db   struct {
		// data source name (ex. database_name, user_name etc...)
		dsn string
	}
	jwt struct {
		secret string
	}
}

// application ... application log & configuration
type application struct {
	config config
	logger *log.Logger
	models models.Models
}

// AppStatus ... application status struct
type AppStatus struct {
	Status      string `json:"status"`
	Environment string `json:"environment"`
	Version     string `json:"version"`
}

func main() {
	var cfg config
	// 設定用のconfigのフィールドをコマンドライン引数から受け取る
	flag.IntVar(&cfg.port, "port", 4000, "Server port to listen on")
	flag.StringVar(&cfg.env, "env", "development", "Application environment (development|production)")
	flag.StringVar(&cfg.db.dsn, "dsn", "postgres://postgres@localhost/manage_movies?sslmode=disable", "Postgres connection starting")
	flag.StringVar(&cfg.jwt.secret, "jwt-secret", "2dce505d96a53c5768052ee90f3df2055657518dad489160df9913f66042e160", "secret")
	flag.Parse()

	// コマンドライン出力用ログを作成する
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	// PostgresDBと接続開始
	db, err := openDB(cfg)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	app := &application{
		config: cfg,
		logger: logger,
		models: models.NewModels(db),
	}

	// APIサーバーを作成
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Println("Starting server on port", cfg.port)

	// 作成したAPIサーバーを立ち上げる
	err = srv.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}

func openDB(cfg config) (*sql.DB, error) {
	// データベースをcfgで指定した情報でオープンする
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	// 5秒後にタイムアウトする
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// DBにPingをつける
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
