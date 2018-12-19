package xdominion

import (
  "fmt"
  "strings"
)

const (
  OP_Equal = "="
  OP_NotEqual = "!="
  OP_Inferior = "<="
  OP_StrictInferior = "<"
  OP_Superior = ">="
  OP_StrictSuperior = ">"
  OP_Between = "between"
  OP_In = "in"
  OP_NotIn = "not in"
  OP_Like = "like"
  OP_NotLike = "not like"
  OP_Match = "match"
  OP_NotMatch = "not match"
)

/*
  The XConditions is a colection of XCondition
*/

type XConditions []XCondition

// =====================
// XConditions
// =====================

func (c *XConditions)CreateConditions(table *XTable, DB string) string {
  cond := ""

fmt.Println(c)

  for _, xc := range *c {
    cond += xc.GetCondition(table, DB)
  }
  return cond
}

/*
  The XCondition structure
*/

type XCondition struct {
  Field string
  Operator string
  Limit interface{}
  LimitExtra interface{}
  OperatorGlobal string
  AtomOpen int
  AtomClose int
}

// =====================
// XCondition
// =====================

func NewXCondition(field string, operator string, limit interface{}, args ...interface{}) XCondition {
  c := XCondition{Field: field, Operator: operator, Limit: limit}
  for i, p := range args {
    if i == 0 { c.OperatorGlobal = p.(string) }
    if i == 1 { c.AtomOpen = p.(int) }
    if i == 2 { c.AtomClose = p.(int) }
    if i == 3 { c.LimitExtra = p }
  }
  return c
}

func (c *XCondition)GetCondition(table *XTable, DB string) string {
  
    field := table.GetField(c.Field);
    
    fmt.Println("Condicion para ", field)
    
    if field == nil {
      return ""
    }
    
    fmt.Println("Creando la condicion: ")

    cond := ""
    
    if len(c.OperatorGlobal) > 0 {
      cond += " " + c.OperatorGlobal + " "
    }
    
    if c.AtomOpen > 0 {
      cond += strings.Repeat("(", c.AtomOpen)
    }
    switch c.Operator {
      case OP_Equal:
        if c.Limit == nil {
          cond += table.Prepend + field.GetName() + " is " + field.CreateValue(nil, table.Name, DB, table.Prepend);
        } else {
          cond += table.Prepend + field.GetName() + OP_Equal + field.CreateValue(c.Limit, table.Name, DB, table.Prepend);
        }
      case OP_NotEqual:
        if c.Limit == nil {
          cond += table.Prepend + field.GetName() + " is not " + field.CreateValue(nil, table.Name, DB, table.Prepend);
        } else {
          cond += table.Prepend + field.GetName() + OP_NotEqual + field.CreateValue(c.Limit, table.Name, DB, table.Prepend);
        }
      case OP_Superior, OP_StrictSuperior, OP_Inferior, OP_StrictInferior:
        cond += table.Prepend + field.GetName() + c.Operator + field.CreateValue(c.Limit, table.Name, DB, table.Prepend);
      case OP_In, OP_NotIn:
        cond += table.Prepend + field.GetName() + c.Operator + fmt.Sprint(c.Limit)
      case OP_Like:
        cond += table.Prepend + field.GetName() + " ilike '%" + fmt.Sprint(c.Limit) + "%'"
      case OP_NotLike:
        cond += table.Prepend + field.GetName() + " not ilike '%" + fmt.Sprint(c.Limit) + "%'"
      case OP_Match:
        cond += table.Prepend + field.GetName() + " ilike '" + fmt.Sprint(c.Limit) + "'"
      case OP_NotMatch:
        cond += table.Prepend + field.GetName() + " not ilike '" + fmt.Sprint(c.Limit) + "'"
    }
    if c.AtomClose > 0 {
      cond += strings.Repeat(")", c.AtomClose)
    }
    
    fmt.Println(cond)
    
    return cond
  }

