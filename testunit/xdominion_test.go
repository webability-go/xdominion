package main

import (
  "fmt"
  "testing"
  "github.com/webability-go/xdominion"
)

func TestBase(t *testing.T) {
  // Test 1: assign a simple parameter string with some comments
  base := &xdominion.Base{
    DBType: xdominion.DB_Postgres,
    Username: "user",
    Password: "pass",
    Database: "test",
    Host: xdominion.DB_Localhost,
    SSL: false,
  }
  base.Logon()
  
  // print what we got
  fmt.Println(base)

}

