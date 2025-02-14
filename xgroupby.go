package xdominion

/*
  The XGroup is an array of XGroupBy structures
*/

type XGroup []XGroupBy

func (g *XGroup) CreateGroup(table *XTable, DB string) string {
	group := ""

	for _, xg := range *g {
		if group != "" {
			group += ", "
		}
		group += xg.GetGroup(table, DB)
	}
	return group
}

/*
  The XGroupBy structure
*/

type XGroupBy struct {
	Field string
}

func (g *XGroupBy) GetGroup(table *XTable, DB string) string {
	// If the field is part of the table, return the field with prepend, else return the field
	f := table.GetField(g.Field)
	if f != nil {
		return table.Prepend + f.GetName()
	}
	return g.Field
}
