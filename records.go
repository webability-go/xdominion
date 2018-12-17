package xdominion

import (
  "fmt"
//  "time"
  
  "github.com/webability-go/xcore"
)

/*
  The XRecords is based on the xcore/XDatasetCollectionDef (basically same funcions) and also a Record Interface
*/

type XRecordsDef interface {
  xcore.XDatasetCollectionDef
}

type XRecords []XRecordDef

// =====================
// XRecords
// =====================

func (r *XRecords)Stringify() string {
  return fmt.Sprint(r)
}

func (r *XRecords)Push(data XRecordDef) {
  *r = append(*r, data)
}

func (r *XRecords)Pop() XRecordDef {
  data := (*r)[len(*r)-1]
  *r = (*r)[:len(*r)-1]
  return data
}

func (r *XRecords)Get(index int) (XRecordDef, bool) {
  if index < 0 || index >= len(*r) { return nil, false }
  return (*r)[index], true
}

func (r *XRecords)Count() int {
  return len(*r)
}

func (r *XRecords)GetData(key string) (interface{}, bool) {
  for i := len(*r)-1; i >= 0; i-- {
    val, ok := (*r)[i].Get(key)
    if ok { return val, true }
  }
  return nil, false
}

func (r *XRecords)GetDataString(key string) (string, bool) {
  v, ok := r.GetData(key)
  if ok { return fmt.Sprint(v), true }
  return "", false
}

func (r *XRecords)GetDataRange(key string) (xcore.XDatasetCollectionDef, bool) {
  v, ok := r.GetData(key)
  // Verify v IS actually a XDatasetCollectionDef
  if ok { return v.(xcore.XDatasetCollectionDef), true }
  return nil, false
}


