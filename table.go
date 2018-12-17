package xdominion

import (
//  "fmt"
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
//  haskey := false
//  hascondition := false
//  hasorder := false
//  haslimit := false
//  hasoffset := false
//  hasfields := false
  for _, p := range args {
    switch p.(type) {
      case int:
      case float64, string, time.Time: // position 0 only
//      case XCondition:
//      case XOrderBy:
//      case XFieldSet:
    }
  }
  
  // 2. creates fields to scan
  sqlf := ""
  item := 0
  fieldslist := []string{}
  for _, f := range t.Fields {
    fname := f.GetName()
//    if hasfields && fname not in Fields { continue }
    if item > 0 { sqlf += ", " }
    sqlf += t.Prepend + fname
    fieldslist = append(fieldslist, fname)
    item++
  }
  if item == 0 { return nil, errors.New("Error: there is no fields to search for") }
  
  sql := "select " + sqlf +" from " + t.Name;
  
  // 3. build condition query
  
  
  // 4. build order by query
  
  
  // 5. exec and dump result
  cursor, err := t.Base.Exec(sql)
  defer cursor.Close()
  if err != nil { return nil, err }
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
  return *onerec, nil
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
      if f.IsNotNull() || f.IsPrimaryKey() { return nil, errors.New("Field " + fname + " is mandatory") 
      } else { continue }
    }

    if f.IsAutoIncrement() && v.(int) == 0 { continue }
    if f.IsPrimaryKey() { primkey = t.Prepend + fname }

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


