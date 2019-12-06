package main

import (
	"fmt"
	"github.com/webability-go/xdominion"
	"testing"
)

func TestQueryPG(t *testing.T) {

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
		DBType:   xdominion.DB_Postgres,
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

func printResult(title string, res interface{}, err error) {
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(title)
	switch res.(type) {
	case *xdominion.XRecord, xdominion.XRecord:
		fmt.Println(res)
	case *xdominion.XRecords, xdominion.XRecords:
		for _, x := range *res.(*xdominion.XRecords) {
			fmt.Println(x)
		}
	}
}

// taken from an exercise from codesignals.com as example
/*
results:
name	country	scored	missed	wins
FC Tokyo	Japan	26	28	1
Fujian	China	24	26	0
Jesus Maria	Argentina	25	23	3
University Blues	Australia	16	25	2
Asse-Lennik	Belgium	27	25	4
Bisamberg W	Austria	25	14	6
Deportivo Moron	Argentina	14	25	9
Lomas Voley	Argentina	23	25	8
Oudegem W	Belgium	26	24	5
UVC Mank W	Austria	21	25	7
*/
func getTableResultsDef(base *xdominion.XBase) *xdominion.XTable {
	// dont forget to have your database in UTF8

	// postgres: createdb --encoding=UTF8 test
	// mysql: CREATE DATABASE mydatabase CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

	t := xdominion.NewXTable("results", "")
	// NEVER call a field "key" => reserved word for many databases
	t.AddField(xdominion.XFieldInteger{Name: "id", Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.PK},
		xdominion.XConstraint{Type: xdominion.AI},
	}}) // ai, pk
	t.AddField(xdominion.XFieldVarChar{Name: "name", Size: 50})
	t.AddField(xdominion.XFieldVarChar{Name: "country", Size: 50})
	t.AddField(xdominion.XFieldInteger{Name: "scored"})
	t.AddField(xdominion.XFieldInteger{Name: "missed"})
	t.AddField(xdominion.XFieldInteger{Name: "wins"})
	t.SetBase(base)

	_, err := base.Exec("drop table results")
	err = t.Synchronize()
	if err != nil {
		fmt.Println(err)
	}

	_, err = t.Insert(xdominion.XRecord{"id": 1, "name": "FC Tokyo", "country": "Japan", "scored": 26, "missed": 28, "wins": 1})
	if err != nil {
		fmt.Println(err)
	}
	// inserts UTF 8
	t.Insert(xdominion.XRecord{"id": 2, "name": "Fujian 研发部门", "country": "China", "scored": 24, "missed": 2, "wins": 0})
	t.Insert(xdominion.XRecord{"id": 3, "name": "Jesus María", "country": "Argentina", "scored": 25, "missed": 23, "wins": 3})
	t.Insert(xdominion.XRecord{"id": 4, "name": "University Blues", "country": "Australia", "scored": 16, "missed": 25, "wins": 2})
	t.Insert(xdominion.XRecord{"id": 5, "name": "Asse-Lennik", "country": "Belgium", "scored": 27, "missed": 25, "wins": 4})
	t.Insert(xdominion.XRecord{"id": 6, "name": "Bisamberg W", "country": "Austria", "scored": 25, "missed": 14, "wins": 6})
	t.Insert(xdominion.XRecord{"id": 7, "name": "Deportivo Moron", "country": "Argentina", "scored": 14, "missed": 25, "wins": 9})
	t.Insert(xdominion.XRecord{"id": 8, "name": "Lomas Voley áñóú", "country": "Argentina", "scored": 23, "missed": 25, "wins": 8})
	t.Insert(xdominion.XRecord{"id": 9, "name": "Oudegem W", "country": "Belgium", "scored": 26, "missed": 24, "wins": 5})
	t.Insert(xdominion.XRecord{"id": 10, "name": "UVC Mank W", "country": "Austria", "scored": 21, "missed": 25, "wins": 7})

	return t
}
