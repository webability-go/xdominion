package xdominion

const (
	PK = "pk"
	NN = "nn"
	AI = "ai"
	FK = "fk"
	IN = "in"
	UI = "ui"
	MI = "mi"
	MU = "mu"
	DC = "dc"
	TR = "tr"
)

/*
  The XConstraints is a colection of XConstrain
*/

type XConstraints []XConstraint

// =====================
// XConstraints
// =====================

func (c *XConstraints) Get(ctype string) *XConstraint {
	// TODO(phil) And what if there are more than one contraint of this type ? for instance MI and MU may be more than one
	for _, ct := range *c {
		if ct.Type == ctype {
			return &ct
		}
	}
	return nil
}

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

/*
  The XConstraint structure
*/

type XConstraint struct {
	Type string
	Data []string
}

// =====================
// XConstraint
// =====================
