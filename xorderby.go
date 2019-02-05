package xdominion

/*
  The XOrderBy is an array of field names
*/

type XOrder []XOrderBy

func (o *XOrder)CreateOrder(table *XTable, DB string) string {
  order := ""

  for _, xo := range *o {
    // , entre cada uno
    order += xo.GetOrder(table, DB)
  }
  return order
}

/*
  The XOrderBy structure
*/

type XOrderBy struct {
  Field string
  Operator string
}

func NewXOrderBy(field string, operator string) XOrderBy {
  o := XOrderBy{Field: field, Operator: operator}
  return o
}

func (c *XOrderBy)GetOrder(table *XTable, DB string) string {
  
    field := table.GetField(c.Field);
    
    if field == nil {
      return ""
    }
    
    order := table.Prepend + field.GetName()
    
    if len(c.Operator) > 0 {
      order += " " + c.Operator
    }
    return order
  }

