package xdominion

import (
  "fmt"
  "time"
  
  "github.com/webability-go/xcore"
)

/*
  The XRecord is based on the xcore/XDataset (basically same funcions) and also a Record Interface
*/

type XRecordDef interface {
  xcore.XDatasetDef
  GetTime(string) (time.Time, bool)
  GetInt(string) (int, bool)
  GetFloat(string) (float64, bool)
}

type XRecord map[string]interface{}

func NewXRecord() *XRecord {
  return &XRecord{}
}


// =====================
// XRecord
// =====================

// makes an interface of XDataset to reuse for otrhe libraries and be sure we can call the functions
func (r *XRecord)Stringify() string {
  return fmt.Sprint(*r)
//  return fmt.Sprintf("DIRECCION DEL OBJETO: %p %p ", &d, d)
}

func (r *XRecord)Get(key string) (interface{}, bool) {
  data, ok := (*r)[key]
  if ok { return data, true }
  return nil, false
}

func (r *XRecord)GetString(key string) (string, bool) {
  data, ok := (*r)[key]
  if ok { return fmt.Sprint(data), true }
  return "", false
}

func (r *XRecord)GetCollection(key string) (xcore.XDatasetCollectionDef, bool) {
  data, ok := (*r)[key]
  if ok { return data.(xcore.XDatasetCollectionDef), true }
  return nil, false
}

func (r *XRecord)Set(key string, data interface{}) {
  (*r)[key] = data
}

func (r *XRecord)Del(key string) {
  delete(*r, key)
}



func (r *XRecord)GetInt(key string) (int, bool) {
  v, ok := (*r)[key]
  if ok {
    return v.(int), true
  }
  return 0, false
}

func (r *XRecord)GetFloat(key string) (float64, bool) {
  v, ok := (*r)[key]
  if ok {
    return v.(float64), true
  }
  return 0.0, false
}

func (r *XRecord)GetTime(key string) (time.Time, bool) {
  v, ok := (*r)[key]
  if ok {
    return v.(time.Time), true
  }
  return time.Time{}, false
}






