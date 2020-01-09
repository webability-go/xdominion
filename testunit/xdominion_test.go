package main

import (
	"fmt"
	"testing"
	//  "github.com/webability-go/xcore"
	"github.com/webability-go/xdominion"
)

func TestBase(t *testing.T) {

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

	// print what we got
	fmt.Println(base)

	fmt.Println("drop table test")
	q, err := base.Exec("drop table test")
	if err != nil {
		fmt.Println(err)
	}
	q.Close()

	// Creates a table
	tb := getTableDef(base)

	err = tb.Synchronize()
	if err != nil {
		fmt.Println(err)
	}

	res1, err := tb.Insert(xdominion.XRecord{"f1": 1, "f2": "Data line 1"})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res1)
	res2, err := tb.Insert(xdominion.XRecord{"f1": 2, "f2": "Data line 2"})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res2)

	res21, err := tb.Insert(xdominion.XRecord{"f1": 2, "f3": "test"})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res21)

	fmt.Println("HACIENDO UN QUERY GENERAL:")
	res3, err := tb.Select()
	if err != nil {
		fmt.Println(err)
	} else {
		for _, x := range *res3.(*xdominion.XRecords) {
			fmt.Println(x)
		}
	}

	res4, err := tb.Select(1)
	if err != nil {
		fmt.Println(err)
	} else {
		switch res4.(type) {
		case *xdominion.XRecord:
			fmt.Println(res4)
		case *xdominion.XRecords:
			for _, x := range *res4.(*xdominion.XRecords) {
				fmt.Println(x)
			}
		}
	}

	res5, err := tb.Select(xdominion.XConditions{xdominion.NewXCondition("f1", "=", 1), xdominion.NewXCondition("f2", "like", "lin", "and")})
	if err != nil {
		fmt.Println(err)
	} else {
		switch res5.(type) {
		case *xdominion.XRecord:
			fmt.Println(res5)
		case *xdominion.XRecords:
			for _, x := range *res5.(*xdominion.XRecords) {
				fmt.Println(x)
			}
		}
	}

	fmt.Println("Ahora vamos a probar funciones que siempre regresan el mismo cast:")
	res6, err := tb.SelectOne(xdominion.XConditions{xdominion.NewXCondition("f1", "=", 1), xdominion.NewXCondition("f1", "=", 2, "or")})
	fmt.Println(res6)
	res7, err := tb.SelectAll(xdominion.XConditions{xdominion.NewXCondition("f1", "=", 1), xdominion.NewXCondition("f1", "=", 2, "and")})
	fmt.Println(res7)
	res8, err := tb.SelectAll(xdominion.XConditions{xdominion.NewXCondition("f1", "=", 1)})
	fmt.Println(res8)

	res9, err := tb.SelectAll(xdominion.XFieldSet{"f1", "f4"})
	fmt.Println(res9)

	i, err := tb.Update(xdominion.XRecord{"f6": 3.1415927})
	fmt.Println("Updated: ", i, err)

	num, err := tb.Count()
	fmt.Println("Count: ", num)

	res10, err := tb.SelectAll(xdominion.XFieldSet{"f1", "f6"})
	for _, r := range *res10 {
		for _, f := range *(r.(*xdominion.XRecord)) {
			fmt.Printf("%T %v: ", f, f)
		}
	}

	res11 := res6.Clone()

	fmt.Println(res11)

}

func getTableDef(base *xdominion.XBase) *xdominion.XTable {
	t := xdominion.NewXTable("test", "t_")
	t.AddField(xdominion.XFieldInteger{Name: "f1", Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.PK},
		xdominion.XConstraint{Type: xdominion.AI},
	}}) // ai, pk
	t.AddField(xdominion.XFieldVarChar{Name: "f2", Size: 20, Constraints: xdominion.XConstraints{
		xdominion.XConstraint{Type: xdominion.NN},
	}})
	t.AddField(xdominion.XFieldText{Name: "f3"})
	t.AddField(xdominion.XFieldDate{Name: "f4"})
	t.AddField(xdominion.XFieldDateTime{Name: "f5"})
	t.AddField(xdominion.XFieldFloat{Name: "f6"})
	t.SetBase(base)
	return t
}

/* Test injection of a recordset into a template */
/*
func TestTemplate(t *testing.T) {

  tmpl, _ := xcore.NewXTemplateFromString(`
Some data:
@@result@@
[[result]]  Data 1: {{f1}}, data 2: {{f2}}
[[]]
End of array of data

`)

  base := &xdominion.XBase{
    DBType: xdominion.DB_Postgres,
    Username: "username",
    Password: "password",
    Database: "test",
    Host: xdominion.DB_Localhost,
    SSL: false,
  }
  base.Logon()

  tb := getTableDef(base)
  irecs, _ := tb.Select()
  recs := irecs.(*xdominion.XRecords)

  // the data must be into "result" parameter
  data := xdominion.NewXRecord()
  data.Set("result", recs)

  fmt.Println(recs)
  fmt.Println(data)

  result := tmpl.Execute(data)
  fmt.Println("Result: ", result)

}

*/
