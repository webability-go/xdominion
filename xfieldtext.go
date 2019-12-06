package xdominion

import (
	"fmt"
)

type XFieldText struct {
	Name        string
	Constraints XConstraints
}

// creates the name of the field with its type (to create the table)
func (f XFieldText) CreateField(prepend string, DB string, ifText *bool) string {
	ftype := " text"
	extra := ""
	if f.Constraints != nil {
		extra = f.Constraints.CreateConstraints(prepend, f.Name, DB)
	}
	return prepend + f.Name + ftype + extra
}

// creates a string representation of the value of the field for insert/update and queries where
func (f XFieldText) CreateValue(v interface{}, table string, DB string, id string) string {
	return "'" + fmt.Sprint(v) + "'"
}

// creates the sequence used by the field (only autoincrement fields)
func (f XFieldText) CreateSequence(table string) string {
	return ""
}

// creates the index used by the field (normal, unique, multi, multi unique)
func (f XFieldText) CreateIndex(table string, id string, DB XBase) string {
	return ""
}

// gets the name of the field
func (f XFieldText) GetName() string {
	return f.Name
}

// gets the type of the field
func (f XFieldText) GetType() int {
	return XField_Text
}

// gets the checks of the field
func (f XFieldText) GetConstraints() XConstraints {
	return f.Constraints
}
