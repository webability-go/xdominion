package xdominion

import (
	"fmt"
	//  "time"
)

type XFieldDateTime struct {
	Name        string
	Constraints XConstraints
}

// creates the name of the field with its type (to create the table)
func (f XFieldDateTime) CreateField(prepend string, DB string, ifText *bool) string {
	ftype := " timestamp"
	extra := ""
	if f.Constraints != nil {
		extra = f.Constraints.CreateConstraints(prepend, f.Name, DB)
	}
	return prepend + f.Name + ftype + extra
}

// creates a string representation of the value of the field for insert/update and queries where
func (f XFieldDateTime) CreateValue(v interface{}, table string, DB string, id string) string {
	return "'" + fmt.Sprint(v) + "'"
}

// creates the sequence used by the field (only autoincrement fields)
func (f XFieldDateTime) CreateSequence(table string) string {
	return ""
}

// creates the index used by the field (normal, unique, multi, multi unique)
func (f XFieldDateTime) CreateIndex(table string, id string, DB XBase) string {
	return ""
}

// gets the name of the field
func (f XFieldDateTime) GetName() string {
	return f.Name
}

// gets the type of the field
func (f XFieldDateTime) GetType() int {
	return XField_DateTime
}

// gets the checks of the field
func (f XFieldDateTime) GetConstraints() XConstraints {
	return f.Constraints
}
