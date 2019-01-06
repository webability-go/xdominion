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
}

type XRecord map[string]interface{}

func NewXRecord() *XRecord {
  return &XRecord{}
}


// =====================
// XRecord
// =====================

// makes an interface of XDataset to reuse for other libraries and be sure we can call the functions
func (r *XRecord)Stringify() string {
  return fmt.Sprint(*r)
}

func (r *XRecord)Set(key string, data interface{}) {
  (*r)[key] = data
}

func (r *XRecord)Get(key string) (interface{}, bool) {
  data, ok := (*r)[key]
  if ok { return data, true }
  return nil, false
}

func (r *XRecord)GetDataset(key string) (xcore.XDatasetDef, bool) {
  if val, ok := (*r)[key]; ok {
    switch val.(type) {
      case *XRecord: return val.(*XRecord), true
    }
  }
  return nil, false
}

func (r *XRecord)GetCollection(key string) (xcore.XDatasetCollectionDef, bool) {
  if val, ok := (*r)[key]; ok {
    switch val.(type) {
      case *XRecords: return val.(*XRecords), true
    }
  }
  return nil, false
}

func (r *XRecord)GetString(key string) (string, bool) {
  if val, ok := (*r)[key]; ok {
    switch val.(type) {
      case string: return val.(string), true
    }
  }
  return "", false
}

func (r *XRecord)GetInt(key string) (int, bool) {
  if val, ok := (*r)[key]; ok {
    switch val.(type) {
      case int: return val.(int), true
    }
  }
  return 0, false
}

func (r *XRecord)GetBool(key string) (bool, bool) {
  if val, ok := (*r)[key]; ok {
    switch val.(type) {
      case bool: return val.(bool), true
    }
  }
  return false, false
}

func (r *XRecord)GetFloat(key string) (float64, bool) {
  if val, ok := (*r)[key]; ok {
    switch val.(type) {
      case float64: return val.(float64), true
    }
  }
  return 0, false
}

func (r *XRecord)GetStringCollection(key string) ([]string, bool) {
  if val, ok := (*r)[key]; ok {
    switch val.(type) {
      case []string: return val.([]string), true
      case string: return []string{val.(string)}, true
    }
  }
  return nil, false
}

func (r *XRecord)GetBoolCollection(key string) ([]bool, bool) {
  if val, ok := (*r)[key]; ok {
    switch val.(type) {
      case []bool: return val.([]bool), true
      case bool: return []bool{val.(bool)}, true
    }
  }
  return nil, false
}

func (r *XRecord)GetIntCollection(key string) ([]int, bool) {
  if val, ok := (*r)[key]; ok {
    switch val.(type) {
      case []int: return val.([]int), true
      case int: return []int{val.(int)}, true
    }
  }
  return nil, false
}

func (r *XRecord)GetFloatCollection(key string) ([]float64, bool) {
  if val, ok := (*r)[key]; ok {
    switch val.(type) {
      case []float64: return val.([]float64), true
      case float64: return []float64{val.(float64)}, true
    }
  }
  return nil, false
}

func (r *XRecord)Del(key string) {
  delete(*r, key)
}




func (r *XRecord)GetTime(key string) (time.Time, bool) {
  v, ok := (*r)[key]
  if ok {
    return v.(time.Time), true
  }
  return time.Time{}, false
}






