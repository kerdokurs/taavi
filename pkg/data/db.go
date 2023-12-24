package data

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var db *sqlx.DB

func Init() {
	db = sqlx.MustConnect("sqlite3", "taavi.db")
}
