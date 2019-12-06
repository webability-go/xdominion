package xdominion

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"log"
	"os"
)

/* IMPORTANT NOTE:
Because of import structure for GO, you need to import the database drivers you will need to connect to.
As of 2018/12/01, only postgres and mysql are supported for now
*/

const (
	// Version of XDominion
	VERSION = "0.1.0"

	// The distinct supported databases
	DB_Postgres  = "postgres"
	DB_MySQL     = "mysql"
	DB_Localhost = "localhost"
)

var DEBUG bool = false

type XBase struct {
	DB       *sql.DB
	Logged   bool
	DBType   string
	Username string
	Password string
	Database string
	Host     string
	SSL      bool
	Logger   *log.Logger
}

func (b *XBase) Logon() {
	if b.Logged {
		return
	}
	b.Logged = true

	if b.Logger == nil {
		b.Logger = log.New(os.Stderr, "XDominion: ", log.LstdFlags)
	}

	var err error
	var src string
	switch b.DBType {
	case DB_Postgres:
		src = b.DBType + "://" + b.Username + ":" + b.Password + "@" + b.Host + "/" + b.Database
		if !b.SSL {
			src += "?sslmode=disable"
		}
	case DB_MySQL:
		if b.Host == DB_Localhost {
			src = b.Username + ":" + b.Password + "@" + "/" + b.Database
		} else {
			src = b.Username + ":" + b.Password + "@" + b.Host + "/" + b.Database
		}
	}

	if DEBUG {
		fmt.Println("DB Source:", src)
	}

	b.DB, err = sql.Open(b.DBType, src)
	if err != nil {
		log.Panic(err)
		return
	}

	if err = b.DB.Ping(); err != nil {
		log.Panic(err)
	}
}

func (b *XBase) Logoff() {

}

func (b *XBase) Exec(query string, args ...interface{}) (*sql.Rows, error) {
	cursor, err := b.DB.Query(query, args...)
	return cursor, err
}

func (b *XBase) Cursor() *Cursor {
	return &Cursor{Base: b, Transactional: false}
}

func (b *XBase) CursorTransactional() *Cursor {
	c := b.Cursor()
	c.BeginTransaction()
	return c
}

/*
func (b *XBase)BeginTransaction() {
  b.DB.Begin()
}

func (b *XBase)Commit() {
  b.DB.Commit()
}

func (b *XBase)Rollback() {
  b.DB.Rollback()
}
*/
