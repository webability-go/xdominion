package main

import (
  "fmt"
  "testing"
  "github.com/webability-go/xdominion"
)

func TestBase(t *testing.T) {
  // Test 1: assign a simple parameter string with some comments
  // Be sure the database exists or have an error
  /*
  // Install postgres server and client on your server
  UNIX> createdb --encoding=UTF-8 test
  UNIX> psql test
  psql> 
  Modifiy pg_hba.conf to authorize ::1 with your user/pass
  */
  base := &xdominion.XBase{
    DBType: xdominion.DB_Postgres,
    Username: "username",
    Password: "password",
    Database: "test",
    Host: xdominion.DB_Localhost,
    SSL: false,
  }
  base.Logon()
  
  // print what we got
  fmt.Println(base)
  
  // Creates a table
  tb := getTableDef(base)

  err := tb.Synchronize()
  if (err != nil) {
    fmt.Println(err)
  }
  
  res, err := tb.Insert(xdominion.XRecord{"f1": 1, "f2": "Data line 1",})
  if (err != nil) {
    fmt.Println(err)
  }
  fmt.Println(res)
  res, err = tb.Insert(xdominion.XRecord{"f1": 2, "f2": "Data line 2",})
  if (err != nil) {
    fmt.Println(err)
  }
  fmt.Println(res)

  fmt.Println("HACIENDO UN QUERY GENERAL:")
  res, err = tb.Select()

  for _, x := range res.(xdominion.XRecords) {
    fmt.Println(x)
  }
  
}

func getTableDef(base *xdominion.XBase) *xdominion.XTable {
  t := xdominion.NewXTable("test", "t_")
  t.AddField(xdominion.XFieldInteger{Name: "f1"})
  t.AddField(xdominion.XFieldVarChar{Name: "f2", Size: 20})
  t.SetBase(base)
  return t
}
