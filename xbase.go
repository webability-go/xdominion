package xdominion

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

/* IMPORTANT NOTE:
Because of import structure for GO, you need to import the database drivers you will need to connect to.
As of 2018/12/01, only postgres and mysql are supported for now
*/

const (
	// The distinct supported databases
	DB_Postgres  = "postgres"
	DB_MySQL     = "mysql"
	DB_Localhost = "localhost"
)

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

func (b *XBase) NewXCursor() *XCursor {
	c := &XCursor{Base: b}
	return c
}

type XTransaction struct {
	DB *XBase
	TX *sql.Tx
}

func (b *XBase) BeginTransaction() (*XTransaction, error) {
	ctx := context.Background()
	tx, err := b.DB.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return nil, err
	}
	return &XTransaction{
		DB: b,
		TX: tx,
	}, nil
}

func (t *XTransaction) Exec(query string, args ...interface{}) (*sql.Rows, error) {
	cursor, err := t.TX.Query(query, args...)
	return cursor, err
}

func (t *XTransaction) NewXCursor() *XCursor {
	c := &XCursor{Transaction: t}
	return c
}

func (t *XTransaction) Commit() error {
	return t.TX.Commit()
}

func (t *XTransaction) Rollback() error {
	return t.TX.Rollback()
}
