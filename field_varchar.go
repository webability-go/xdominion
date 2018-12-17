package xdominion

import (
  "fmt"
)

type XFieldVarChar struct
{
  Name string
  Size int
}

// creates the name of the field with its type (to create the table)
func (f XFieldVarChar)CreateField(prepend string, DB string, ifText *bool) string {
  return prepend + f.Name + " varchar(" + fmt.Sprint(f.Size) + ")"
}

// creates a string representation of the value of the field for insert/update
func (f XFieldVarChar)CreateValue(v interface{}, table string, DB string, id string) string {
  return fmt.Sprint(v)
}

// creates the sequence used by the field (only autoincrement fields)
func (f XFieldVarChar)CreateSequence(table string) string {
  return ""
}

// creates the index used by the field (normal, unique, multi, multi unique)
func (f XFieldVarChar)CreateIndex(table string, id string, DB XBase) string {
  return ""
}

// gets the name of the field
func (f XFieldVarChar)GetName() string {
  return f.Name
}

// gets the type of the field
func (f XFieldVarChar)GetType() int {
  return XField_VarChar
}

// gets the checks of the field
func (f XFieldVarChar)GetChecks() interface{} {
  return nil
}

// returns true if the field is a primary key for the table
func (f XFieldVarChar)IsPrimaryKey() bool {
  return false
}

// returns true if the field is an auto-incremented field (with a sequence)
func (f XFieldVarChar)IsAutoIncrement() bool {
  return false
}

// returns true if the field cannot be null
func (f XFieldVarChar)IsNotNull() bool {
  return false
}

// returns true if the field checks contains a specific condition
func (f XFieldVarChar)Contains(check string) bool {
  return false
}

// returns the foreign key of the field if defined
func (f XFieldVarChar)GetForeignKey() string {
  return ""
}
