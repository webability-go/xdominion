package xdominion

/*
  The XGroup is an array of XGroupBy structures
*/

type XGroup []XGroupBy

func (g *XGroup)CreateGroup(table *XTable, DB string) string {
  group := ""

  for _, xg := range *g {
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

func (g *XGroupBy)GetGroup(table *XTable, DB string) string {
  return "Group By --"
}

