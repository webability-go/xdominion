package xdominion

//  "fmt"

const (
	XField_Int      = 1
	XField_VarChar  = 2
	XField_Float    = 3
	XField_DateTime = 4
	XField_Date     = 5
	XField_Text     = 6
)

// XFieldDef is an interface representing a field definition.
// It provides methods for creating the field, its value,
// its sequence, and its index. It also provides methods for
// getting the field's name, type, and constraints.
type XFieldDef interface {
	// CreateField creates the name of the field with its type (to create the table).
	// The `prepend` argument is used to add a prefix to the field name, and
	// the `DB` argument is used to specify the database type (postgres or mysql).
	// The `ifText` argument is a boolean pointer that is used to track whether
	// the field is a text field or not.
	CreateField(prepend string, DB string, ifText *bool) string

	// CreateValue creates a string representation of the value of the field for
	// insert/update with ' for text. The `v` argument is the value of the field,
	// the `table` argument is the name of the table, the `DB` argument is used to
	// specify the database type (postgres or mysql), and the `id` argument is used
	// to specify the ID of the row (for updates).
	CreateValue(v interface{}, table string, DB string, id string) string

	// CreateSequence creates the sequence used by the field (only autoincrement fields).
	// The `table` argument is the name of the table.
	CreateSequence(table string) string

	// CreateIndex creates the index used by the field (normal, unique, multi, multi unique).
	// The `table` argument is the name of the table, the `id` argument is the ID of the row,
	// and the `DB` argument is used to specify the database type (postgres or mysql).
	CreateIndex(table string, id string, DB string) []string

	// GetName gets the name of the field.
	GetName() string

	// GetType gets the type of the field.
	GetType() int

	// GetConstraints gets the checks of the field.
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

// IsPrimaryKey returns true if the field is a primary key for the table.
func IsPrimaryKey(f XFieldDef) bool {
	return IsFieldConstraint(f, PK)
}

// IsNotNull returns true if the field cannot be null.
func IsNotNull(f XFieldDef) bool {
	return IsFieldConstraint(f, NN)
}

// IsAutoIncrement returns true if the field is an auto-incremented field (with a sequence).
func IsAutoIncrement(f XFieldDef) bool {
	return IsFieldConstraint(f, AI)
}

// IsFieldConstraint returns true if the field checks contains a specific condition.
func IsFieldConstraint(f XFieldDef, ftype string) bool {
	xc := f.GetConstraints()
	if xc == nil {
		return false
	}
	for _, c := range xc {
		if c.Type == ftype {
			return true
		}
	}
	return false
}
