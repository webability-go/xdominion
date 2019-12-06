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
