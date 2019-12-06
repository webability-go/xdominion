package main

import (
	"github.com/webability-go/xdominion"
	"testing"
)

func TestQueryMY(t *testing.T) {

	xdominion.DEBUG = true
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
		DBType:   xdominion.DB_MySQL,
		Username: "username",
		Password: "password",
		Database: "test",
		Host:     xdominion.DB_Localhost,
		SSL:      false,
	}
	base.Logon()

	// Creates a table
	tb := getTableResultsDef(base)

	// Queries
	res, err := tb.Select()
	printResult("The whole table, raw:", res, err)

	// Queries simple
	res, err = tb.Select(xdominion.XOrder{xdominion.XOrderBy{Field: "scored", Operator: xdominion.ASC}, xdominion.XOrderBy{Field: "missed", Operator: xdominion.DESC}})
	printResult("The whole table, ordered:", res, err)

	// wheres
	res, err = tb.Select(xdominion.NewXCondition("country", "=", "Argentina"), xdominion.XOrderBy{Field: "wins", Operator: xdominion.DESC})
	printResult("Only Argentina:", res, err)

	// start with Aust
	res, err = tb.Select(xdominion.NewXCondition("country", "like", "Aust%"), xdominion.XOrderBy{Field: "wins", Operator: xdominion.DESC})
	printResult("Start with Aust:", res, err)

	// contains u insentisive
	res, err = tb.Select(xdominion.NewXCondition("name", "ilike", "%u%"), xdominion.XOrderBy{Field: "wins", Operator: xdominion.DESC})
	printResult("Contains u:", res, err)

}
