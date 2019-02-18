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
  OP_iLike = "ilike"
  OP_NotiLike = "not ilike"
)

/*
  The XConditions is a colection of XCondition
*/

type XConditions []XCondition

// =====================
// XConditions
// =====================

func (c *XConditions)CreateConditions(table *XTable, DB string, baseindex int) (string, []string) {
  cond := ""
  data := []string{}

  for _, xc := range *c {
    scond, sdata, indexused := xc.GetCondition(table, DB, baseindex)
    cond += scond
    data = append(data, sdata)
    if indexused {
      baseindex ++
    }
  }
  return cond, data
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

func (c *XCondition)GetCondition(table *XTable, DB string, baseindex int) (string, string, bool) {
  
    field := table.GetField(c.Field);
    
    if field == nil {
      return "", "", false
    }
    
    cond := ""
    
    if len(c.OperatorGlobal) > 0 {
      cond += " " + c.OperatorGlobal + " "
    }
    
    if c.AtomOpen > 0 {
      cond += strings.Repeat("(", c.AtomOpen)
    }
    indexused := true
    value := ""
    switch c.Operator {
      case OP_Equal:
        if c.Limit == nil {
          cond += table.Prepend + field.GetName() + " is null";
          indexused = false
        } else {
          value = field.GetValue(c.Limit, table.Name, DB, table.Prepend)
          cond += table.Prepend + field.GetName() + OP_Equal + getQueryString(DB, baseindex)
        }
      case OP_NotEqual:
        if c.Limit == nil {
          cond += table.Prepend + field.GetName() + " is not null";
          indexused = false
        } else {
          value = field.GetValue(c.Limit, table.Name, DB, table.Prepend)
          cond += table.Prepend + field.GetName() + OP_NotEqual + getQueryString(DB, baseindex)
        }
      case OP_Superior, OP_StrictSuperior, OP_Inferior, OP_StrictInferior:
        value = field.GetValue(c.Limit, table.Name, DB, table.Prepend)
        cond += table.Prepend + field.GetName() + c.Operator + getQueryString(DB, baseindex)
      case OP_In, OP_NotIn:
        value = fmt.Sprint(c.Limit)
        cond += table.Prepend + field.GetName() + c.Operator + getQueryString(DB, baseindex)
      case OP_Like:
        value = fmt.Sprint(c.Limit)
        cond += table.Prepend + field.GetName() + " like " + getQueryString(DB, baseindex)
      case OP_NotLike:
        value = fmt.Sprint(c.Limit)
        cond += table.Prepend + field.GetName() + " not like " + getQueryString(DB, baseindex)
      case OP_iLike:
        value = fmt.Sprint(c.Limit)
        switch DB {
          case DB_Postgres:
            cond += table.Prepend + field.GetName() + " ilike " + getQueryString(DB, baseindex)
          case DB_MySQL:
            cond += "lower(" + table.Prepend + field.GetName() + ") like lower(" + getQueryString(DB, baseindex) + ")"
        }
      case OP_NotiLike:
        value = fmt.Sprint(c.Limit)
        switch DB {
          case DB_Postgres:
            cond += table.Prepend + field.GetName() + " not ilike " + getQueryString(DB, baseindex)
          case DB_MySQL:
            cond += "lower(" + table.Prepend + field.GetName() + ") not like lower(" + getQueryString(DB, baseindex) + ")"
        }
      default:
        // warning: operator not supported
    }
    if c.AtomClose > 0 {
      cond += strings.Repeat(")", c.AtomClose)
    }
    
    return cond, value, indexused
  }

