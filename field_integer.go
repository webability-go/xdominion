package xdominion

import (
  "fmt"
)

type XFieldInteger struct
{
  Name string
}

// creates the name of the field with its type (to create the table)
func (f XFieldInteger)CreateField(prepend string, DB string, ifText *bool) string {
  return prepend + f.Name + " integer"
}

// creates a string representation of the value of the field for insert/update
func (f XFieldInteger)CreateValue(v interface{}, table string, DB string, id string) string {
  return fmt.Sprint(v)
}

// creates the sequence used by the field (only autoincrement fields)
func (f XFieldInteger)CreateSequence(table string) string {
  return ""
}

// creates the index used by the field (normal, unique, multi, multi unique)
func (f XFieldInteger)CreateIndex(table string, id string, DB XBase) string {
  return ""
}

// gets the name of the field
func (f XFieldInteger)GetName() string {
  return f.Name
}

// gets the type of the field
func (f XFieldInteger)GetType() int {
  return XField_Int
}

// gets the checks of the field
func (f XFieldInteger)GetChecks() interface{} {
  return nil
}

// returns true if the field is a primary key for the table
func (f XFieldInteger)IsPrimaryKey() bool {
  return false
}

// returns true if the field is an auto-incremented field (with a sequence)
func (f XFieldInteger)IsAutoIncrement() bool {
  return false
}

// returns true if the field cannot be null
func (f XFieldInteger)IsNotNull() bool {
  return false
}

// returns true if the field checks contains a specific condition
func (f XFieldInteger)Contains(check string) bool {
  return false
}

// returns the foreign key of the field if defined
func (f XFieldInteger)GetForeignKey() string {
  return ""
}
