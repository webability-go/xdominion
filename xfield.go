package xdominion

import (
//  "fmt"
)

const (
  XField_Int        = 1
  XField_VarChar    = 2
  XField_Float      = 3
  XField_DateTime   = 4
  XField_Date       = 5
  XField_Text       = 6
)

type XFieldDef interface {
  // creates the name of the field with its type (to create the table)
  CreateField(prepend string, DB string, ifText *bool) string
  // creates a string representation of the value of the field for insert/update with ' for text
  CreateValue(v interface{}, table string, DB string, id string) string
  // creates the sequence used by the field (only autoincrement fields)
  CreateSequence(table string) string
  // creates the index used by the field (normal, unique, multi, multi unique)
  CreateIndex(table string, id string, DB XBase) string
  // gets the name of the field
  GetName() string
  // gets the type of the field
  GetType() int
  // gets the checks of the field
  GetConstraints() XConstraints
  // returns true if the field is a primary key for the table
//  IsPrimaryKey() bool
  // returns true if the field is an auto-incremented field (with a sequence)
//  IsAutoIncrement() bool
  // returns true if the field cannot be null
//  IsNotNull() bool
  // returns true if the field checks contains a specific condition
//  Contains(check string) bool
  // returns the foreign key of the field if defined
//  GetForeignKey() string
}

// returns true if the field is a primary key for the table
func IsPrimaryKey(f XFieldDef) bool {
  return IsFieldConstraint(f, PK)
}

func IsNotNull(f XFieldDef) bool {
  return IsFieldConstraint(f, NN)
}

func IsAutoIncrement(f XFieldDef) bool {
  return IsFieldConstraint(f, AI)
}

func IsFieldConstraint(f XFieldDef, ftype string) bool {
  xc := f.GetConstraints()
  if xc == nil { return false }
  for _, c := range xc {
    if c.Type == ftype { return true }
  }
  return false
}

