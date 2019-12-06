package xdominion

import (
	"fmt"
)

type XFieldVarChar struct {
	Name        string
	Size        int
	Constraints XConstraints
}

// creates the name of the field with its type (to create the table)
func (f XFieldVarChar) CreateField(prepend string, DB string, ifText *bool) string {
	ftype := " varchar(" + fmt.Sprint(f.Size) + ")"
	extra := ""
	if f.Constraints != nil {
		extra = f.Constraints.CreateConstraints(prepend, f.Name, DB)
	}
	return prepend + f.Name + ftype + extra
}

// creates a string representation of the value of the field for insert/update and queries where
func (f XFieldVarChar) CreateValue(v interface{}, table string, DB string, id string) string {
	return "'" + fmt.Sprint(v) + "'"
}

// creates the sequence used by the field (only autoincrement fields)
func (f XFieldVarChar) CreateSequence(table string) string {
	return ""
}

// creates the index used by the field (normal, unique, multi, multi unique)
func (f XFieldVarChar) CreateIndex(table string, id string, DB XBase) string {
	return ""
}

// gets the name of the field
func (f XFieldVarChar) GetName() string {
	return f.Name
}

// gets the type of the field
func (f XFieldVarChar) GetType() int {
	return XField_VarChar
}

// gets the checks of the field
func (f XFieldVarChar) GetConstraints() XConstraints {
	return f.Constraints
}
