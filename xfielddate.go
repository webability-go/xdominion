package xdominion

import (
  "fmt"
)

type XFieldDate struct
{
  Name string
  Constraints XConstraints
}

// creates the name of the field with its type (to create the table)
func (f XFieldDate)CreateField(prepend string, DB string, ifText *bool) string {
  ftype := " date"
  extra := ""
  if f.Constraints != nil {
    extra = f.Constraints.CreateConstraints(prepend, f.Name, DB)
  }
  return prepend + f.Name + ftype + extra
}

// creates a string representation of the value of the field for insert/update and queries where
func (f XFieldDate)CreateValue(v interface{}, table string, DB string, id string) string {
  return "'" + fmt.Sprint(v) + "'"
}

// gets directly the value of the field for insert/update and queries where
func (f XFieldDate)GetValue(v interface{}, table string, DB string, id string) string {
  return fmt.Sprint(v)
}

// creates the sequence used by the field (only autoincrement fields)
func (f XFieldDate)CreateSequence(table string) string {
  return ""
}

// creates the index used by the field (normal, unique, multi, multi unique)
func (f XFieldDate)CreateIndex(table string, id string, DB XBase) string {
  return ""
}

// gets the name of the field
func (f XFieldDate)GetName() string {
  return f.Name
}

// gets the type of the field
func (f XFieldDate)GetType() int {
  return XField_Date
}

// gets the checks of the field
func (f XFieldDate)GetConstraints() XConstraints {
  return f.Constraints
}
