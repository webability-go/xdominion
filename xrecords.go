package xdominion

import (
	"fmt"
	"time"

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

func (r *XRecords) String() string {
	return r.GoString()
}

func (r *XRecords) GoString() string {
	return fmt.Sprint(r)
}

func (d *XRecords) Unshift(data xcore.XDatasetDef) {
	*d = append([]XRecordDef{data.(*XRecord)}, (*d)...)
}

func (d *XRecords) Shift() xcore.XDatasetDef {
	data := (*d)[0]
	*d = (*d)[1:]
	return data
}

func (r *XRecords) Push(data xcore.XDatasetDef) {
	*r = append(*r, data.(*XRecord))
}

func (r *XRecords) Pop() xcore.XDatasetDef {
	data := (*r)[len(*r)-1]
	*r = (*r)[:len(*r)-1]
	return data
}

func (r *XRecords) Count() int {
	return len(*r)
}

func (r *XRecords) Get(index int) (xcore.XDatasetDef, bool) {
	if index < 0 || index >= len(*r) {
		return nil, false
	}
	return (*r)[index], true
}

func (r *XRecords) GetData(key string) (interface{}, bool) {
	for i := len(*r) - 1; i >= 0; i-- {
		val, ok := (*r)[i].Get(key)
		if ok {
			return val, true
		}
	}
	return nil, false
}

func (r *XRecords) GetDataString(key string) (string, bool) {
	v, ok := r.GetData(key)
	if ok {
		return fmt.Sprint(v), true
	}
	return "", false
}

func (d *XRecords) GetDataBool(key string) (bool, bool) {
	if val, ok := d.GetData(key); ok {
		switch val.(type) {
		case bool:
			return val.(bool), true
		}
	}
	return false, false
}

func (d *XRecords) GetDataInt(key string) (int, bool) {
	if val, ok := d.GetData(key); ok {
		switch val.(type) {
		case int:
			return val.(int), true
		}
	}
	return 0, false
}

func (d *XRecords) GetDataFloat(key string) (float64, bool) {
	if val, ok := d.GetData(key); ok {
		switch val.(type) {
		case float64:
			return val.(float64), true
		}
	}
	return 0, false
}

func (d *XRecords) GetDataTime(key string) (time.Time, bool) {
	if val, ok := d.GetData(key); ok {
		switch val.(type) {
		case time.Time:
			return val.(time.Time), true
		}
	}
	return time.Time{}, false
}

func (r *XRecords) GetCollection(key string) (xcore.XDatasetCollectionDef, bool) {
	if val, ok := r.GetData(key); ok {
		switch val.(type) {
		case *XRecords:
			return val.(*XRecords), true
		}
	}
	return nil, false
}

func (r *XRecords) Clone() xcore.XDatasetCollectionDef {
	cloned := &XRecords{}
	for _, val := range *r {
		*cloned = append(*cloned, val.Clone())
	}
	return cloned
}
