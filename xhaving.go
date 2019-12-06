package xdominion

/*
  The XOrderBy is an array of field names
*/

type XHaving []XCondition

func (h *XHaving) CreateHaving(table *XTable, DB string, baseindex int) (string, []interface{}) {
	having := ""
	data := []interface{}{}

	for _, xh := range *h {
		shaving, sdata, indexused := xh.GetCondition(table, DB, baseindex)
		having += shaving
		data = append(data, sdata)
		if indexused {
			baseindex++
		}
	}
	return having, data
}
