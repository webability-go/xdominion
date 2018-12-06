package xdominion

import (

)

type Table struct {

}

func NewTable() *Table {
  
}

func NewTableDef(file string) *Table {
  
}

func (t *Table)AddDef(name string, prepend string) {
  
}

func (t *Table)AddField(field *Field) {
  
}

func (t *Table)DoSelect(key interface{}, fields *[]string) *Record {
  
}

func (t *Table)DoSelectCondition(condition *Condition, orderby *OrderBy, first int, last int, fields *[]string) *Records {
  
}

func (t *Table)DoUpdate(key interface{}, record *Record) {
  
}

func (t *Table)DoUpdateCondition(condition *Condition, record *Record) {
  
}

func (t *Table)DoDelete(key interface{}) {
  
}

func (t *Table)DoDeleteCondition(condition *Condition) {
  
}

func (t *Table)DoInsert(record *Record) {
  
}





