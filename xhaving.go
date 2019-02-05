package xdominion

/*
  The XOrderBy is an array of field names
*/

type XHaving []XCondition

func (h *XHaving)CreateHaving(table *XTable, DB string) string {
  having := ""

  for _, xh := range *h {
    having += xh.GetCondition(table, DB)
  }
  return having
}

