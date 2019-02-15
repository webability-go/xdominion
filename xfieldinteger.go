package xdominion

import (
  "fmt"
)

type XFieldInteger struct
{
  Name string
  Constraints XConstraints
}

// creates the name of the field with its type (to create the table)
func (f XFieldInteger)CreateField(prepend string, DB string, ifText *bool) string {
  ftype := " integer"
  if IsAutoIncrement(f) {
    ftype = " serial"
  }
  extra := ""
  if f.Constraints != nil {
    extra = f.Constraints.CreateConstraints(prepend, f.Name, DB)
  }
  return prepend + f.Name + ftype + extra
}

// creates a string representation of the value of the field for insert/update
func (f XFieldInteger)CreateValue(v interface{}, table string, DB string, id string) string {
  return fmt.Sprint(v)
}

// gets directly the value of the field for insert/update and queries where
func (f XFieldInteger)GetValue(v interface{}, table string, DB string, id string) string {
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
func (f XFieldInteger)GetConstraints() XConstraints {
  return f.Constraints
}
