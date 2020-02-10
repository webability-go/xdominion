package xdominion

import (
	"fmt"
	"strconv"
	"time"

	"github.com/webability-go/xcore"
)

/*
  The XRecord is based on the xcore/XDataset (basically same funcions) and also a Record Interface
*/

type XRecordDef interface {
	xcore.XDatasetDef
}

type XRecord map[string]interface{}

func NewXRecord() *XRecord {
	return &XRecord{}
}

// =====================
// XRecord
// =====================

// makes an interface of XDataset to reuse for other libraries and be sure we can call the functions
func (r *XRecord) String() string {
	return r.GoString()
}

func (r *XRecord) GoString() string {
	return fmt.Sprint(*r)
}

func (r *XRecord) Set(key string, data interface{}) {
	(*r)[key] = data
}

func (r *XRecord) Get(key string) (interface{}, bool) {
	data, ok := (*r)[key]
	if ok {
		return data, true
	}
	return nil, false
}

func (r *XRecord) GetDataset(key string) (xcore.XDatasetDef, bool) {
	if val, ok := (*r)[key]; ok {
		switch val.(type) {
		case *XRecord:
			return val.(*XRecord), true
		}
	}
	return nil, false
}

func (r *XRecord) GetCollection(key string) (xcore.XDatasetCollectionDef, bool) {
	if val, ok := (*r)[key]; ok {
		switch val.(type) {
		case *XRecords:
			return val.(*XRecords), true
		}
	}
	return nil, false
}

func (r *XRecord) GetString(key string) (string, bool) {
	if val, ok := (*r)[key]; ok {
		switch val.(type) {
		case nil:
			return "", false
		case string:
			return val.(string), true
		default:
			return fmt.Sprint(val), true
		}
	}
	return "", false
}

func (r *XRecord) GetInt(key string) (int, bool) {
	if val, ok := (*r)[key]; ok {
		switch val.(type) {
		case int:
			return val.(int), true
		case int64:
			return int(val.(int64)), true
		case float64:
			return int(val.(float64)), true
		case bool:
			if val.(bool) {
				return 1, true
			} else {
				return 0, true
			}
		}
	}
	return 0, false
}

func (r *XRecord) GetBool(key string) (bool, bool) {
	if val, ok := (*r)[key]; ok {
		switch val.(type) {
		case bool:
			return val.(bool), true
		case int:
			return val.(int) != 0, true
		case int64:
			return val.(int64) != 0, true
		case float64:
			return val.(float64) != 0, true
		}
	}
	return false, false
}

func (r *XRecord) GetFloat(key string) (float64, bool) {
	if val, ok := (*r)[key]; ok {
		switch val.(type) {
		case string:
			v, err := strconv.ParseFloat(val.(string), 64)
			return v, err == nil
		case float32:
			return float64(val.(float32)), true
		case float64:
			return val.(float64), true
		case int:
			return float64(val.(int)), true
		case int64:
			return float64(val.(int64)), true
		case bool:
			if val.(bool) {
				return 1.0, true
			} else {
				return 0.0, true
			}
		}
	}
	return 0, false
}

func (r *XRecord) GetTime(key string) (time.Time, bool) {
	v, ok := (*r)[key]
	if ok {
		if v != nil {
			return v.(time.Time), true
		}
	}
	return time.Time{}, false
}

func (r *XRecord) GetStringCollection(key string) ([]string, bool) {
	if val, ok := (*r)[key]; ok {
		switch val.(type) {
		case []string:
			return val.([]string), true
		case string:
			return []string{val.(string)}, true
		}
	}
	return nil, false
}

func (r *XRecord) GetBoolCollection(key string) ([]bool, bool) {
	if val, ok := (*r)[key]; ok {
		switch val.(type) {
		case []bool:
			return val.([]bool), true
		case bool:
			return []bool{val.(bool)}, true
		}
	}
	return nil, false
}

func (r *XRecord) GetIntCollection(key string) ([]int, bool) {
	if val, ok := (*r)[key]; ok {
		switch val.(type) {
		case []int:
			return val.([]int), true
		case int:
			return []int{val.(int)}, true
		}
	}
	return nil, false
}

func (r *XRecord) GetFloatCollection(key string) ([]float64, bool) {
	if val, ok := (*r)[key]; ok {
		switch val.(type) {
		case []float64:
			return val.([]float64), true
		case float64:
			return []float64{val.(float64)}, true
		}
	}
	return nil, false
}

func (r *XRecord) GetTimeCollection(key string) ([]time.Time, bool) {
	if val, ok := (*r)[key]; ok {
		switch val.(type) {
		case []time.Time:
			return val.([]time.Time), true
		case time.Time:
			return []time.Time{val.(time.Time)}, true
		}
	}
	return nil, false
}

func (r *XRecord) Del(key string) {
	delete(*r, key)
}

func (r *XRecord) Clone() xcore.XDatasetDef {
	cloned := &XRecord{}
	for id, val := range *r {
		clonedval := val
		// If the object is also cloneable, we clone it
		if cloneable, ok := val.(interface{ Clone() xcore.XDatasetDef }); ok {
			clonedval = cloneable.Clone()
		}
		cloned.Set(id, clonedval)
	}
	return cloned
}
