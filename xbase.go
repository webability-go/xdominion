package xdominion

import (
//  "fmt"
  "log"
  "database/sql"
  _ "github.com/lib/pq"
)

/* IMPORTANT NOTE:
Because of import structure for GO, you need to import the database drivers you will need to connect to.
As of 2018/12/01, only postgres is supported for now
*/


const (
  // Version of XDominion
  VERSION = "0.0.7"

  // The distinct supported databases
  DB_Postgres = "postgres"
  DB_Localhost = "localhost"
)

type XBase struct {
  DB *sql.DB
  Logged bool
  DBType string
  Username string
  Password string
  Database string
  Host string
  SSL bool
}

func (b *XBase)Logon() {
  if b.Logged {
    return
  }
  b.Logged = true

  var err error
  src := b.DBType + "://" + b.Username + ":" + b.Password + "@" + b.Host + "/" + b.Database
  if !b.SSL {
    src += "?sslmode=disable"
  }

  b.DB, err = sql.Open("postgres", src)
  if err != nil {
    log.Panic(err)
  }

  if err = b.DB.Ping(); err != nil {
    log.Panic(err)
  }
}

func (b *XBase)Logoff() {

}

func (b *XBase)Exec(query string, args ...interface{}) (*sql.Rows, error) {
  cursor, err := b.DB.Query(query, args...)
  return cursor, err
}

func (b *XBase)Cursor() *Cursor {
  return &Cursor{Base: b, Transactional: false,}
}

func (b *XBase)CursorTransactional() *Cursor {
  c := b.Cursor()
  c.BeginTransaction()
  return c
}

