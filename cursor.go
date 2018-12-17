package xdominion

import (
//  "fmt"
  "log"
  "errors"
//  "database/sql"
  _ "github.com/lib/pq"
)

type Cursor struct {
  Base *XBase
  Transactional bool
}

func (c *Cursor)Query() interface{} {
  
  return nil
}

func (c *Cursor)Close() {
  if c.Transactional {
    c.Commit()
  }
}

func (c *Cursor)BeginTransaction() {
  if c.Transactional {
    log.Panic(errors.New("A transation has already started in cursor"))
  }
  c.Transactional = true
}

func (c *Cursor)Rollback() {
  if !c.Transactional {
    log.Panic(errors.New("There is no transaction declared in cursor"))
  }

  c.Transactional = false
}

func (c *Cursor)Commit() {
  if !c.Transactional {
    log.Panic(errors.New("There is no transaction declared in cursor"))
  }

  c.Transactional = false
}
