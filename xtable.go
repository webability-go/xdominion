package xdominion

import (
  "fmt"
  "strconv"
  "errors"
  "time"
)

type XTable struct {
  Base *XBase
  Name string
  Prepend string
  Fields []XFieldDef
  InsertedKey interface{}
}

func NewXTable(name string, prepend string) *XTable {
  return &XTable{Name: name, Prepend: prepend,}
}

func (t *XTable)AddField(field XFieldDef) {
  t.Fields = append(t.Fields, field)
}

func (t *XTable)SetBase(base *XBase) {
  t.Base = base
}

func (t *XTable)Synchronize() error {
  // This funcion is supposed to check structure vs def and insert/modifie table fiels and constraints
  var ifText bool = false
  
  // creates the "create table" query
  query := "Create table " + t.Name + " ("
  for i, f := range t.Fields {
    if i > 0 { query += "," }
    query += f.CreateField(t.Prepend, t.Base.DBType, &ifText)
  }
  query += ")"
  
  fmt.Println(query)
  
  cursor, err := t.Base.Exec(query)
  if err != nil { return err }
  defer cursor.Close()
  return nil
}

/*
  Select:
  Args are:
  NO ARGS: select * from table
  1rst ARG is a simple cast (int, string, time, float) => primary key. IGNORE other args
  1rst ARG is a XCondition: select where XCondition, then:
    2nd ARG is XOrderby: apply orderby
    3rd ARG is int: limit
    4th ARG is int: offset
    5th ARG is []string: list of fields to get back
  returns nil, error / XRecord, nil / XRecords, nil
*/

func (t *XTable)Select(args ...interface{}) (interface{}, error) {
  // 1. analyse params
  haskey := false
  var key interface{}
  hasconditions := false
  var conditions XConditions
  hasorder := false
  var order XOrderBy
  haslimit := false
  var limit int
  hasoffset := false
  var offset int
  hasfields := false
  var fields XFieldSet
  
  for i, p := range args {
    fmt.Printf("TYPE PARAM: %T\n", p)
    switch p.(type) {
      case int:
        if i == 0 { 
          haskey = true
          key = p
        } else {
          if haslimit {
            hasoffset = true
            offset = p.(int)
          } else {
            haslimit = true
            limit = p.(int)
          }
        }
      case float64, string, time.Time: // position 0 only
        if i == 0 { 
          haskey = true
          key = p
        } else {
          return nil, errors.New("Error: a key value can only be on first parameter")
        }
      case XCondition:
        hasconditions = true
        conditions = XConditions{p.(XCondition)}
      case XConditions:
        hasconditions = true
        conditions = p.(XConditions)
      case XOrderBy:
        hasorder = true
        order = p.(XOrderBy)
      case XFieldSet:
        hasfields = true
        fields = p.(XFieldSet)
    }
  }

fmt.Println(hasconditions)
fmt.Println(conditions)
  
  // 2. creates fields to scan
  sqlf := ""
  item := 0
  fieldslist := []string{}
  for _, f := range t.Fields {
    fname := f.GetName()
    if hasfields {
      fmt.Println(fields)
//    && fname not in Fields { continue }
    }
    if item > 0 { sqlf += ", " }
    sqlf += t.Prepend + fname
    fieldslist = append(fieldslist, fname)
    item++
  }
  if item == 0 { return nil, errors.New("Error: there is no fields to search for") }
  
  sql := "select " + sqlf +" from " + t.Name;
  
  itemdata := 1
  sqldata := make([]interface{}, 0)
  // 3. build condition query
  if haskey {
    // get primary key field
    primkey := t.GetPrimaryKey()
    if primkey == nil {
      return nil, errors.New("There is no primary key on in the table")
    }
    sql += " where " + t.Prepend + primkey.GetName() + " = $" + strconv.Itoa(itemdata)
    sqldata = append(sqldata, primkey.CreateValue(key, t.Name, t.Base.DBType, t.Prepend))
    itemdata++
  } else if hasconditions {
    sql += " where " + conditions.CreateConditions(t, t.Base.DBType)
  }
  
  // 4. build order by query
  if hasorder {
    sql += fmt.Sprint(order)
  }
  
  // 5. Limits
  if haslimit {
    sql += " limit " + strconv.Itoa(limit)
  }
  if hasoffset {
    sql += " offset " + strconv.Itoa(offset)
  }
  
  fmt.Println(sql)

      // 6. exec and dump result
  cursor, err := t.Base.Exec(sql, sqldata...)
  if err != nil { return nil, err }
  defer cursor.Close()

  var result = make([]interface{}, item)
  var bridge = make([]interface{}, item)
  for i, _ := range result {
    bridge[i] = &result[i] // Put pointers to each string in the interface slice
  }
  
  var onerec *XRecord = nil
  var somerecs *XRecords = nil
  for cursor.Next() {
    // scan into result through bridge
    cursor.Scan(bridge...)
    // creates a XRecord with result
    xr := &XRecord{}
    for i, f := range fieldslist {
      xr.Set(f, result[i])
    }
    
    if somerecs != nil {
      somerecs.Push(xr)
    } else {
      if onerec == nil {
        onerec = xr
      } else {
        somerecs = &XRecords{onerec, xr}
        onerec = nil
      }
    }
  }
  if somerecs != nil {
    return *somerecs, nil
  }
  if onerec != nil {
    return *onerec, nil
  }
  return nil, nil
}

// Insert things into database:
// If data is XRecord, insert the record. Returns the key (same type as field type) interface{}
// If data is XRecords, insert the collection of XRecord. Returns an array of keys (same type as field type) []interface{}
// If data is SubQuery, intert the result of the sub query, return ?
func (t *XTable)Insert(data interface{}) (interface{}, error) {

  t.InsertedKey = nil
  switch data.(type) {
    case XRecords:
      rc := data.(XRecords)
      nrc := len(rc)
      keys := make([]interface{}, nrc)
      for _, record := range rc {
         k, err := t.Insert(record);
         if (err != nil) { return nil, err }
         keys = append(keys, k)
      }
      t.InsertedKey = keys
      return keys, nil
    case XRecord:
    default:

// SUB QUERY:
//      $sql = "insert into ".$this->name." ".$record->getSubQuery();
    
      return nil, errors.New("Type of Data no known. Must be one of XRecord, XRecords, SubQuery")
  }

  rc := data.(XRecord)
  sqlf := ""
  sqlv := ""
  item := 0

  primkey := ""
  sqldata := make([]interface{}, 0)
  for _, f := range t.Fields {
    fname := f.GetName()
    v, ok := rc.Get(fname)
    if !ok { 
      if IsNotNull(f) || IsPrimaryKey(f) { return nil, errors.New("Field " + fname + " is mandatory") 
      } else { continue }
    }

    if IsAutoIncrement(f) && v.(int) == 0 { continue }
    if IsPrimaryKey(f) { primkey = t.Prepend + fname }

    if item > 0 {
      sqlf += ", "
      sqlv += ", "
    }
    sqlf += t.Prepend + fname
    sqldata = append(sqldata, f.CreateValue(v, t.Name, t.Base.DBType, t.Prepend))
    sqlv += "$" + strconv.Itoa(item+1)
    item++
  }
  sql := "insert into " + t.Name + " ("+sqlf+") values ("+sqlv+")";

  if primkey != "" {
    sql += " returning " + primkey;
    var key interface{}
    cursor, err := t.Base.Exec(sql, sqldata...)
    if err != nil { return nil, err }
    defer cursor.Close()
    cursor.Next()
    err = cursor.Scan(&key)
    if err != nil { return nil, err }
    t.InsertedKey = key
    return key, nil
  }

  cursor, err := t.Base.Exec(sql, sqldata...)
  if err != nil { return nil, err }
  defer cursor.Close()
  return nil, nil
}

func (t *XTable)Update(key interface{}, record *XRecord) error {
  return nil
}

func (t *XTable)UpdateCondition(record *XRecord, condition interface{}) error {
  return nil
}

func (t *XTable)Delete(key interface{}) error {
  return nil
}

func (t *XTable)DeleteCondition(condition interface{}) error {
  return nil
}

func (t *XTable)Count(condition interface{}) (int, error) {
  return 0, nil
}

func (t *XTable)Min(field string, condition interface{}) (interface{}, error) {
  return nil, nil
}

func (t *XTable)Max(field string, condition interface{}) (interface{}, error) {
  return nil, nil
}

func (t *XTable)Avg(field string, condition interface{}) (interface{}, error) {
  return nil, nil
}

func (t *XTable)GetPrimaryKey() XFieldDef {
  for _, f := range t.Fields {
    if IsPrimaryKey(f) { return f }
  }
  return nil
}

func (t *XTable)GetField(name string) XFieldDef {
  for _, f := range t.Fields {
    if f.GetName() == name { return f }
  }
  return nil
}
