package xdominion

const (
	PK = "pk" // Primary Key Constraint type
	NN = "nn" // Not Null Constraint type
	AI = "ai" // Auto Increment Constraint type
	FK = "fk" // Foreign Key Constraint type
	IN = "in" // Index Constraint type
	UI = "ui" // Unique Index Constraint type
	MI = "mi" // Multiple Index Constraint type
	MU = "mu" // Multiple Unique Index Constraint type
	DC = "dc" // Drop cascade Constraint type
	TR = "tr" // Transfer to another PK before deleting (is not drop cascade)
)

// XConstraints is a collection of XConstraint structures.
type XConstraints []XConstraint

// =====================
// XConstraints
// =====================

// Get function returns a pointer to an XConstraint structure whose type matches the ctype string.
func (c *XConstraints) Get(ctype string) *XConstraint {
	// TODO(phil) And what if there are more than one contraint of this type ? for instance MI and MU may be more than one
	for _, ct := range *c {
		if ct.Type == ctype {
			return &ct
		}
	}
	return nil
}

// CreateConstraints function returns a string of constraints for a database field.
// prepend: a string to be prepended before the field name.
// name: the name of the field.
// DB: the database type.
func (c *XConstraints) CreateConstraints(prepend string, name string, DB string) string {
	cnt := ""
	pk := c.Get(PK)
	if pk != nil {
		cnt += " primary key"
	}
	ai := c.Get(AI)
	if ai != nil {
		if DB == DB_MySQL {
			cnt += " auto_increment"
		}
		// mysql, mssql: build auto_increment and indentity(1,1). pgsql ya tiene un tipo "serial" desde el campo mismo
	}
	nn := c.Get(NN)
	if nn != nil {
		cnt += " not null"
	}
	fk := c.Get(FK)
	if fk != nil {
		cnt += " references " + fk.Data[0] + "(" + fk.Data[1] + ")"
	}
	return cnt
}

// CreateIndex function returns an array of strings of SQL index creation queries.
// table: the name of the table.
// prepend: a string to be prepended before the field name.
// field: the name of the field.
// DB: the database type.
func (c *XConstraints) CreateIndex(table string, prepend string, field string, DB string) []string {

	// TODO(phil) simplify the code, it's virtually the same code for 4 types of indexes
	indexes := []string{}
	in := c.Get(IN)
	if in != nil {
		i := "create index i" + prepend + field + " on " + table + "(" + prepend + field + ")"
		indexes = append(indexes, i)
	}
	ui := c.Get(UI)
	if ui != nil {
		i := "create unique index i" + prepend + field + " on " + table + "(" + prepend + field + ")"
		indexes = append(indexes, i)
	}
	mi := c.Get(MI)
	if mi != nil {
		flds := prepend + field
		for _, f := range mi.Data {
			flds += "," + prepend + f
		}
		i := "create index i" + prepend + field + " on " + table + "(" + flds + ")"
		indexes = append(indexes, i)
	}
	mu := c.Get(MU)
	if mu != nil {
		flds := prepend + field
		for _, f := range mu.Data {
			flds += "," + prepend + f
		}
		i := "create unique index i" + prepend + field + " on " + table + "(" + flds + ")"
		indexes = append(indexes, i)
	}
	return indexes
}

// XConstraint is a structure representing a database constraint.
// Type: the type of constraint, one of PK, NN, AI, FK, IN, UI, MI, MU, DC and TR.
// Data: an array of strings containing the data related to the constraint.
type XConstraint struct {
	Type string
	Data []string
}
