package xdominion

import (
	"fmt"
	"testing"
)

func ExampleXBase() {

	base := &XBase{
		DBType:   DB_Postgres,
		Username: "username",
		Password: "password",
		Database: "test",
		Host:     DB_Localhost,
		SSL:      false,
	}
	base.Logon()
}

func TestXCore_commit(t *testing.T) {

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

	fmt.Println(tb.Synchronize())

	// insert
	res1, err := tb.Insert(XRecord{"f1": 3, "f2": "Data line 1"})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res1)
	res2, err := tb.Insert(XRecord{"f1": 4, "f2": "Data line 2"})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res2)

	tx, err := base.BeginTransaction()

	// insert
	res3, err := tb.Insert(XRecord{"f1": 5, "f2": "Data line 1"}, tx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res3)

	res4, err := tb.Insert(XRecord{"f1": 6, "f2": "Data line 2"}, tx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res4)

	res41, err := tb.Update(XRecord{"f3": "some new data"}, tx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res41)

	res42, err := tb.Delete(1, tx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res42)

	// lets roll back
	err = tx.Commit()
	if err != nil {
		fmt.Println(err)
	}

	// should give an error, tx is rollbacked
	res5, err := tb.Insert(XRecord{"f1": 7, "f2": "Data line 2"}, tx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res5)

}

func getTableDef(base *XBase) *XTable {
	t := NewXTable("test", "t_")
	t.AddField(XFieldInteger{Name: "f1", Constraints: XConstraints{
		XConstraint{Type: PK},
		XConstraint{Type: AI},
	}}) // ai, pk
	t.AddField(XFieldVarChar{Name: "f2", Size: 20, Constraints: XConstraints{
		XConstraint{Type: NN},
	}})
	t.AddField(XFieldText{Name: "f3"})
	t.AddField(XFieldDate{Name: "f4", Constraints: XConstraints{XConstraint{Type: IN}}})
	t.AddField(XFieldDateTime{Name: "f5", Constraints: XConstraints{XConstraint{Type: UI}}})
	t.AddField(XFieldFloat{Name: "f6", Constraints: XConstraints{XConstraint{Type: MU, Data: []string{"f4", "f5"}}}})
	t.SetBase(base)
	return t
}
