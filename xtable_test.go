package xdominion

import (
	"fmt"
	"testing"
)

func TestXTable_select(t *testing.T) {

	base := &XBase{
		DBType:   DB_Postgres,
		Username: "username",
		Password: "password",
		Database: "test",
		Host:     DB_Localhost,
		SSL:      false,
	}
	base.Logon()

	tb := getTableDef(base)
	tb.Synchronize()
	buildData(tb)

	var where *XCondition
	var order *XOrderBy

	r, err := tb.SelectAll(where, order)
	fmt.Println(r, err)

	r, err = tb.SelectAll(&XFieldSet{"f2", "count(*)"}, &XConditions{NewXCondition("f1", "!=", 1)}, &XGroupBy{"f2"}, &XOrderBy{"count(*)", DESC}, 10)
	fmt.Println(r, err)
}

func buildData(tb *XTable) {
	// insert some things
	_, err := tb.Insert(XRecord{"f1": 3, "f2": "Data line 1"})
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = tb.Insert(XRecord{"f1": 4, "f2": "Data line 2"})
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = tb.Insert(XRecord{"f1": 5, "f2": "Data line 1"})
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = tb.Insert(XRecord{"f1": 6, "f2": "Data line 2"})
	if err != nil {
		fmt.Println(err)
		return
	}
}
