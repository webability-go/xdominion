package xdominion

import (
	"fmt"
	"testing"
)

func TestXCursor_exec(t *testing.T) {

	base := &XBase{
		DBType:   DB_Postgres,
		Username: "username",
		Password: "password",
		Database: "test",
		Host:     DB_Localhost,
		SSL:      false,
	}
	base.Logon()

	fmt.Println("XCURSOR:")

	// 1. Create a cursor
	cursor := base.NewXCursor()
	err := cursor.Exec("select * from test")
	fmt.Println(err, cursor)

	for cursor.Next() {
		data, err := cursor.Read()
		fmt.Println(err, data)
	}
	cursor.Close()

}
