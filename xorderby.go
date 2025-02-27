package xdominion

/*
  The XOrderBy is an array of field names
*/

const (
	ASC  = "asc"
	DESC = "desc"
)

type XOrder []XOrderBy

func (o *XOrder) CreateOrder(table *XTable, DB string) string {
	order := ""

	item := 0
	for _, xo := range *o {
		// , entre cada uno
		if item > 0 {
			order += ", "
		}
		order += xo.GetOrder(table, DB)
		item++
	}
	return order
}

func (o *XOrder) Clone() XOrder {
	no := XOrder{}
	for _, xo := range *o {
		nxo := xo.Clone()
		no = append(no, nxo)
	}
	return no
}

/*
  The XOrderBy structure
*/

type XOrderBy struct {
	Field    string
	Operator string
}

func NewXOrderBy(field string, operator string) XOrderBy {
	o := XOrderBy{Field: field, Operator: operator}
	return o
}

func (c *XOrderBy) GetOrder(table *XTable, DB string) string {

	field := table.GetField(c.Field)

	if field == nil {
		return c.Field + " " + c.Operator
	}

	order := table.Prepend + field.GetName()

	if len(c.Operator) > 0 {
		order += " " + c.Operator
	}
	return order
}

func (o *XOrderBy) Clone() XOrderBy {
	no := XOrderBy{
		Field:    o.Field,
		Operator: o.Operator,
	}
	return no
}
