package main

import (
  	"time"
    "database/sql"
    "github.com/uptrace/bun"
    "github.com/uptrace/bun/dialect/pgdialect"
    "github.com/uptrace/bun/driver/pgdriver"
)

type Note struct {
	bun.BaseModel `bun:"note"`
	Id            int       `json:"id" bun:",scanonly"`
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	UserID        int       `json:"userid" bun:"userid"`
	LastUpdated   time.Time `json:"last_updated" bun:"default:current_timestamp"`
}

type User struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	FullName string `json:"fullname"`
	Email    string `json:"email"`
}


func loadDBCursor() {
	dsn := "postgres://notetoself:notetoself@localhost:5432/notetoself?sslmode=disable"
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	DB = bun.NewDB(sqldb, pgdialect.New())
}
