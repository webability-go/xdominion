package xdominion

import (
	libsql "database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	"golang.org/x/text/language"
)

type XTable struct {
	Base    *XBase
	Name    string
	Prepend string
	// content of table language, informative
	Language    language.Tag
	Fields      []XFieldDef
	InsertedKey interface{}
}

func NewXTable(name string, prepend string) *XTable {
	return &XTable{Name: name, Prepend: prepend}
}

func (t *XTable) AddField(field XFieldDef) {
	t.Fields = append(t.Fields, field)
}

func (t *XTable) SetBase(base *XBase) {
	t.Base = base
}

func (t *XTable) SetLanguage(lang language.Tag) {
	t.Language = lang
}

func (t *XTable) Synchronize(args ...interface{}) error {
	// This funcion is supposed to check structure vs def and insert/modifie table fiels and constraints
	var ifText bool = false
	hastrx := false
	var trx *XTransaction

	if len(args) > 0 {
		trx, hastrx = args[0].(*XTransaction)
		if trx == nil {
			hastrx = false
		}
	}

	// creates the "create table" query
	sql := "Create table " + t.Name + " ("
	indexes := []string{}
	for i, f := range t.Fields {
		if i > 0 {
			sql += ","
		}
		sql += f.CreateField(t.Prepend, t.Base.DBType, &ifText)
		indexes = append(indexes, f.CreateIndex(t.Name, t.Prepend, t.Base.DBType)...)
	}
	sql += ")"

	if DEBUG {
		fmt.Println(sql)
	}

	var cursor *libsql.Rows
	var err error
	if hastrx {
		cursor, err = trx.Exec(sql)
	} else {
		cursor, err = t.Base.Exec(sql)
	}
	if err != nil {
		return err
	}
	cursor.Close()

	if len(indexes) > 0 {
		for _, sqli := range indexes {
			if hastrx {
				cursor, err = trx.Exec(sqli)
			} else {
				cursor, err = t.Base.Exec(sqli)
			}
			if err != nil {
				return err
			}
			cursor.Close()
		}
	}
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
    6th ARG is bool: true = returns always one record max and no more (force limit = 1) and return an XRecord always
  returns nil, error / XRecord, nil / XRecords, nil; or nil, error / XRecords, nil
*/

func (t *XTable) Select(args ...interface{}) (interface{}, error) {
	// 1. analyse params
	haskey := false
	var key interface{}
	hasconditions := false
	var conditions XConditions
	hasorder := false
	var order XOrder
	haslimit := false
	var limit int
	hasoffset := false
	var offset int
	hasgroup := false
	var group XGroup
	hasfields := false
	var fields XFieldSet
	onlyone := false
	var ok bool
	hastrx := false
	var trx *XTransaction

	for i, p := range args {
		switch p.(type) {
		case bool:
			if p.(bool) {
				onlyone = true
			}
		case int, int32, int64:
			if i == 0 {
				haskey = true
				key = p
			} else {
				if haslimit {
					hasoffset = true
					offset, ok = p.(int)
					if !ok {
						offset32, ok := p.(int32)
						if !ok {
							offset64 := p.(int64)
							offset = int(offset64)
						} else {
							offset = int(offset32)
						}
					}
				} else {
					haslimit = true
					limit, ok = p.(int)
					if !ok {
						limit32, ok := p.(int32)
						if !ok {
							limit64 := p.(int64)
							limit = int(limit64)
						} else {
							limit = int(limit32)
						}
					}
				}
			}
		case float32, float64, string, time.Time: // position 0 only
			if i == 0 {
				haskey = true
				key = p
			} else {
				return nil, errors.New("Error: a key value can only be on first parameter")
			}
		case XCondition:
			hasconditions = true
			conditions = XConditions{p.(XCondition)}
		case *XCondition:
			if p != nil && p.(*XCondition) != nil {
				fmt.Printf("DENTRO DEL TEST: %v, %p\n", p, p)

				hasconditions = true
				conditions = XConditions{*p.(*XCondition)}
			}
		case XConditions:
			hasconditions = true
			conditions = p.(XConditions)
		case *XConditions:
			if p != nil && p.(*XConditions) != nil {
				hasconditions = true
				conditions = *p.(*XConditions)
			}
		case XOrderBy:
			hasorder = true
			order = XOrder{p.(XOrderBy)}
		case *XOrderBy:
			if p != nil && p.(*XOrderBy) != nil {
				hasorder = true
				order = XOrder{*p.(*XOrderBy)}
			}
		case XOrder:
			hasorder = true
			order = p.(XOrder)
		case *XOrder:
			if p != nil && p.(*XOrder) != nil {
				hasorder = true
				order = *p.(*XOrder)
			}
		case XGroupBy:
			hasgroup = true
			group = XGroup{p.(XGroupBy)}
		case *XGroupBy:
			if p != nil && p.(*XGroupBy) != nil {
				hasgroup = true
				group = XGroup{*p.(*XGroupBy)}
			}
		case XGroup:
			hasgroup = true
			group = p.(XGroup)
		case *XGroup:
			if p != nil && p.(*XGroup) != nil {
				hasgroup = true
				group = *p.(*XGroup)
			}
		case XFieldSet:
			hasfields = true
			fields = p.(XFieldSet)
		case *XFieldSet:
			if p != nil && p.(*XFieldSet) != nil {
				hasfields = true
				fields = *p.(*XFieldSet)
			}
		case *XTransaction:
			trx, hastrx = p.(*XTransaction)
			if trx == nil {
				hastrx = false
			}
		}
	}
	if onlyone {
		limit = 1
		haslimit = true
	}

	// 2. creates fields to scan
	sqlf := ""
	item := 0
	fieldslist := []string{}
	for _, f := range t.Fields {
		fname := f.GetName()
		if hasfields {
			infield := false
			for _, fn := range fields {
				if fn == fname {
					infield = true
					break
				}
			}
			if !infield {
				continue
			}
		}
		if item > 0 {
			sqlf += ", "
		}
		sqlf += t.Prepend + fname
		fieldslist = append(fieldslist, fname)
		item++
	}
	if item == 0 {
		return nil, errors.New("Error: there is no fields to search for")
	}

	sql := "select " + sqlf + " from " + t.Name

	itemdata := 1
	sqldata := make([]interface{}, 0)
	// 3. build condition query
	if haskey {
		// get primary key field
		primkey := t.GetPrimaryKey()
		if primkey == nil {
			return nil, errors.New("There is no primary key on in the table")
		}
		sql += " where " + t.Prepend + primkey.GetName() + " = " + getQueryString(t.Base.DBType, itemdata)
		sqldata = append(sqldata, key)
		itemdata++
	} else if hasconditions {
		scond, vars := conditions.CreateConditions(t, t.Base.DBType, itemdata)
		sql += " where " + scond
		for _, v := range vars {
			sqldata = append(sqldata, v)
		}
	}
	if DEBUG {
		fmt.Println(sqldata)
	}

	// group by, needs a fieldset
	if hasgroup {
		sql += " group by " + group.CreateGroup(t, t.Base.DBType)
	}

	// 4. build order by query
	if hasorder {
		sql += " order by " + order.CreateOrder(t, t.Base.DBType)
	}

	// having, needs a group by, set of conditions

	// 5. Limits
	if haslimit && limit > 0 {
		sql += " limit " + strconv.Itoa(limit)
	}
	if hasoffset {
		sql += " offset " + strconv.Itoa(offset)
	}

	if DEBUG {
		fmt.Println(sql)
	}

	// 6. exec and dump result
	var cursor *libsql.Rows
	var err error
	if hastrx {
		cursor, err = trx.Exec(sql, sqldata...)
	} else {
		cursor, err = t.Base.Exec(sql, sqldata...)
	}
	if err != nil {
		return nil, err
	}
	defer cursor.Close()

	var result = make([]interface{}, item)
	var bridge = make([]interface{}, item)
	for i := range result {
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
			switch result[i].(type) {
			case []byte: // special case returned by mysql for string :S
				xr.Set(f, string(result[i].([]byte)))
			default:
				xr.Set(f, result[i])
			}
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
		return somerecs, nil
	}
	if onerec != nil {
		return onerec, nil
	}
	return nil, nil
}

func (t *XTable) SelectOne(args ...interface{}) (*XRecord, error) {
	args = append(args, true)
	r, err := t.Select(args...) // select only one
	if r == nil || err != nil {
		return nil, err
	}
	switch r.(type) {
	case *XRecord:
		return r.(*XRecord), nil
	case *XRecords:
		onerec, _ := r.(*XRecords).Get(0)
		return onerec.(*XRecord), nil
	}
	return nil, nil
}

func (t *XTable) SelectAll(args ...interface{}) (*XRecords, error) {
	r, err := t.Select(args...) // select all
	if err != nil {
		return nil, err
	}
	if r == nil {
		return &XRecords{}, nil
	}
	switch r.(type) {
	case *XRecord:
		return &XRecords{r.(*XRecord)}, nil
	case *XRecords:
		return r.(*XRecords), nil
	}
	return nil, nil
}

// Insert things into database:
// If data is XRecord, insert the record. Returns the key (same type as field type) interface{}
// If data is XRecords, insert the collection of XRecord. Returns an array of keys (same type as field type) []interface{}
// If data is SubQuery, intert the result of the sub query, return ?
func (t *XTable) Insert(data interface{}, args ...interface{}) (interface{}, error) {

	hastrx := false
	var trx *XTransaction

	t.InsertedKey = nil
	switch data.(type) {
	case XRecords:
		xdata := data.(XRecords)
		return t.Insert(&xdata, args...)
	case *XRecords, XRecordsDef:
		rc := data.(XRecordsDef)
		nrc := rc.Count()
		keys := make([]interface{}, nrc)
		for i := 0; i < nrc; i++ {
			xt, _ := rc.Get(i)
			k, err := t.Insert(xt, args...)
			if err != nil {
				return nil, err
			}
			keys = append(keys, k)
		}
		t.InsertedKey = keys
		return keys, nil
	case XRecord:
		xdata := data.(XRecord)
		return t.Insert(&xdata, args...)
	case *XRecord, XRecordDef:
	default:

		// SUB QUERY:
		//      $sql = "insert into ".$this->name." ".$record->getSubQuery();

		return nil, errors.New("Type of Data no known. Must be one of XRecord, XRecords, SubQuery")
	}

	if len(args) > 0 {
		trx, hastrx = args[0].(*XTransaction)
		if trx == nil {
			hastrx = false
		}
	}

	rc := data.(XRecordDef)
	sqlf := ""
	sqlv := ""
	item := 0

	primkey := ""
	sqldata := make([]interface{}, 0)
	for _, f := range t.Fields {
		fname := f.GetName()
		v, ok := rc.Get(fname)
		if !ok {
			if IsNotNull(f) || IsPrimaryKey(f) {
				return nil, errors.New("Field " + fname + " is mandatory")
			} else {
				continue
			}
		}

		if IsPrimaryKey(f) {
			primkey = t.Prepend + fname
		}
		if IsAutoIncrement(f) && v.(int) == 0 {
			continue
		}

		if item > 0 {
			sqlf += ", "
			sqlv += ", "
		}
		sqlf += t.Prepend + fname
		sqldata = append(sqldata, v)
		sqlv += getQueryString(t.Base.DBType, item+1)
		item++
	}
	sql := "insert into " + t.Name + " (" + sqlf + ") values (" + sqlv + ")"

	var cursor *libsql.Rows
	var err error

	if primkey != "" {
		if t.Base.DBType == DB_Postgres {
			sql += " returning " + primkey
		}

		if DEBUG {
			fmt.Println(sql)
			fmt.Println(sqldata)
		}

		var key interface{}
		if hastrx {
			cursor, err = trx.Exec(sql, sqldata...)
		} else {
			cursor, err = t.Base.Exec(sql, sqldata...)
		}
		if err != nil {
			return nil, err
		}
		defer cursor.Close()

		if t.Base.DBType == DB_Postgres {
			cursor.Next()
			err = cursor.Scan(&key)
			if err != nil {
				return nil, err
			}
			t.InsertedKey = key
		}
		return key, nil
	}

	if DEBUG {
		fmt.Println(sql)
		fmt.Println(sqldata)
	}

	if hastrx {
		cursor, err = trx.Exec(sql, sqldata...)
	} else {
		cursor, err = t.Base.Exec(sql, sqldata...)
	}
	if err != nil {
		return nil, err
	}
	defer cursor.Close()
	return nil, nil
}

/*
  Update:
  Args are:
  NO ARGS: error
  1rst ARG is a simple cast (int, string, time, float) => primary key.
  1rst ARG is a XCondition: select where XCondition, then:
  2nd ARG is XRecord to modify
  If first arg does not exists, the update is applied to the whole table (aka first arg is XRecord)

  returns int, error. int is the quantity of modified records (always 1 if primary key)
*/

func (t *XTable) Update(args ...interface{}) (int, error) {
	// 1. analyse params
	haskey := false
	var key interface{}
	hasconditions := false
	var conditions XConditions
	hasrecord := false
	var record XRecordDef
	hastrx := false
	var trx *XTransaction

	for _, p := range args {
		switch p.(type) {
		case int, int32, int64, float32, float64, string, time.Time: // position 0 only
			haskey = true
			key = p
		case *XTransaction:
			trx, hastrx = p.(*XTransaction)
			if trx == nil {
				hastrx = false
			}
		case XCondition:
			hasconditions = true
			conditions = XConditions{p.(XCondition)}
		case *XCondition:
			if p != nil && p.(*XCondition) != nil {
				hasconditions = true
				conditions = XConditions{*p.(*XCondition)}
			}
		case XConditions:
			hasconditions = true
			conditions = p.(XConditions)
		case *XConditions:
			if p != nil && p.(*XConditions) != nil {
				hasconditions = true
				conditions = *p.(*XConditions)
			}
		case XRecord:
			hasrecord = true
			np := p.(XRecord)
			var xp interface{} = &np
			record = xp.(XRecordDef)
		case *XRecord, XRecordDef:
			hasrecord = true
			record = p.(XRecordDef)
		}
	}
	if !hasrecord {
		return 0, errors.New("Error: there is no record data to use to modify the records of the table")
	}

	item := 0

	itemdata := 1
	sqldata := make([]interface{}, 0)
	sqlf := ""

	for _, f := range t.Fields {
		fname := f.GetName()
		fd, ok := record.Get(fname)
		if !ok {
			continue
		}

		if item > 0 {
			sqlf += ", "
		}
		sqlf += t.Prepend + fname + " = " + getQueryString(t.Base.DBType, itemdata)
		sqldata = append(sqldata, fd)
		item++
		itemdata++
	}
	if item == 0 {
		return 0, errors.New("Error: there is no fields to update")
	}

	sql := "update " + t.Name + " set " + sqlf

	if haskey {
		// get primary key field
		primkey := t.GetPrimaryKey()
		if primkey == nil {
			return 0, errors.New("There is no primary key on in the table")
		}
		sql += " where " + t.Prepend + primkey.GetName() + " = " + getQueryString(t.Base.DBType, itemdata)
		sqldata = append(sqldata, key)
		itemdata++
	} else if hasconditions {
		scond, vars := conditions.CreateConditions(t, t.Base.DBType, itemdata)
		sql += " where " + scond
		for _, v := range vars {
			sqldata = append(sqldata, v)
		}
	}

	if DEBUG {
		fmt.Println(sql)
	}

	var cursor *libsql.Rows
	var err error
	if hastrx {
		cursor, err = trx.Exec(sql, sqldata...)
	} else {
		cursor, err = t.Base.Exec(sql, sqldata...)
	}
	if err != nil {
		return 0, err
	}
	defer cursor.Close()

	return 1, nil
}

/*
  Upsert: update or insert
  Args are:
  NO ARGS: error
  1rst ARG is a simple cast (int, string, time, float) => primary key.
  1rst ARG is a XCondition: select where XCondition and check if exists, then:
  2nd ARG is XRecord to modify
*/

func (t *XTable) Upsert(args ...interface{}) (int, error) {
	// 1. analyse params
	haskey := false
	var key interface{}
	hasconditions := false
	var conditions XConditions
	hasrecord := false
	var record XRecordDef
	hastrx := false
	var trx *XTransaction

	for _, p := range args {
		switch p.(type) {
		case int, float64, string, time.Time: // position 0 only
			haskey = true
			key = p
		case *XTransaction:
			trx, hastrx = p.(*XTransaction)
			if trx == nil {
				hastrx = false
			}
		case XCondition:
			hasconditions = true
			conditions = XConditions{p.(XCondition)}
		case *XCondition:
			if p != nil && p.(*XCondition) != nil {
				hasconditions = true
				conditions = XConditions{*p.(*XCondition)}
			}
		case XConditions:
			hasconditions = true
			conditions = p.(XConditions)
		case *XConditions:
			if p != nil && p.(*XConditions) != nil {
				hasconditions = true
				conditions = *p.(*XConditions)
			}
		case XRecord:
			hasrecord = true
			np := p.(XRecord)
			var xp interface{} = &np
			record = xp.(XRecordDef)
		case *XRecord, XRecordDef:
			hasrecord = true
			record = p.(XRecordDef)
		}
	}
	if !hasrecord {
		return 0, errors.New("Error: there is no record data to use to insert or modify the records of the table")
	}

	// search record
	var rec *XRecord
	if haskey {
		rec, _ = t.SelectOne(key)
	} else if hasconditions {
		rec, _ = t.SelectOne(conditions)
	}
	primkey := t.GetPrimaryKey()
	if rec != nil {
		thekey, _ := rec.Get(primkey.GetName())
		if hastrx {
			return t.Update(thekey, record, trx)
		} else {
			return t.Update(thekey, record)
		}
	}
	hasprimkey, _ := record.Get(primkey.GetName())
	if hasprimkey == nil {
		record.Set(primkey.GetName(), 0)
	}
	var e error
	if hastrx {
		_, e = t.Insert(record, trx)
	} else {
		_, e = t.Insert(record)
	}
	if e != nil {
		return 0, e
	}
	return 1, nil
}

func (t *XTable) Delete(args ...interface{}) (int, error) {
	// 1. analyse params
	haskey := false
	var key interface{}
	hasconditions := false
	var conditions XConditions
	hastrx := false
	var trx *XTransaction

	for _, p := range args {
		switch p.(type) {
		case int, int32, int64, float32, float64, string, time.Time:
			haskey = true
			key = p
		case *XTransaction:
			trx, hastrx = p.(*XTransaction)
			if trx == nil {
				hastrx = false
			}
		case XCondition:
			hasconditions = true
			conditions = XConditions{p.(XCondition)}
		case *XCondition:
			if p != nil && p.(*XCondition) != nil {
				hasconditions = true
				conditions = XConditions{*p.(*XCondition)}
			}
		case XConditions:
			hasconditions = true
			conditions = p.(XConditions)
		case *XConditions:
			if p != nil && p.(*XConditions) != nil {
				hasconditions = true
				conditions = *p.(*XConditions)
			}
		}
	}

	itemdata := 1
	sqldata := make([]interface{}, 0)

	sql := "delete from " + t.Name

	if haskey {
		// get primary key field
		primkey := t.GetPrimaryKey()
		if primkey == nil {
			return 0, errors.New("There is no primary key on in the table")
		}
		sql += " where " + t.Prepend + primkey.GetName() + " = " + getQueryString(t.Base.DBType, itemdata)
		sqldata = append(sqldata, key)
		itemdata++
	} else if hasconditions {
		scond, vars := conditions.CreateConditions(t, t.Base.DBType, itemdata)
		sql += " where " + scond
		for _, v := range vars {
			sqldata = append(sqldata, v)
		}
	}

	if DEBUG {
		fmt.Println(sql)
	}

	var cursor *libsql.Rows
	var err error
	if hastrx {
		cursor, err = trx.Exec(sql, sqldata...)
	} else {
		cursor, err = t.Base.Exec(sql, sqldata...)
	}
	if err != nil {
		return 0, err
	}
	defer cursor.Close()

	return 1, nil
}

func (t *XTable) Count(args ...interface{}) (int, error) {
	// 1. analyse params
	hasfield := false
	var field string
	hasconditions := false
	var conditions XConditions
	hastrx := false
	var trx *XTransaction

	for _, p := range args {
		switch p.(type) {
		case string:
			hasfield = true
			field = p.(string)
		case XCondition:
			hasconditions = true
			conditions = XConditions{p.(XCondition)}
		case *XCondition:
			if p != nil && p.(*XCondition) != nil {
				hasconditions = true
				conditions = XConditions{*p.(*XCondition)}
			}
		case XConditions:
			hasconditions = true
			conditions = p.(XConditions)
		case *XConditions:
			if p != nil && p.(*XConditions) != nil {
				hasconditions = true
				conditions = *p.(*XConditions)
			}
		case *XTransaction:
			trx, hastrx = p.(*XTransaction)
			if trx == nil {
				hastrx = false
			}
		}
	}

	itemdata := 1
	sqldata := make([]interface{}, 0)

	sql := "select count("
	if hasfield {
		sql += "distinct " + t.Prepend + field
	} else {
		sql += "*"
	}
	sql += ") from " + t.Name

	if hasconditions {
		scond, vars := conditions.CreateConditions(t, t.Base.DBType, itemdata)
		sql += " where " + scond
		for _, v := range vars {
			sqldata = append(sqldata, v)
		}
	}

	if DEBUG {
		fmt.Println(sql)
	}

	var cursor *libsql.Rows
	var err error
	if hastrx {
		cursor, err = trx.Exec(sql, sqldata...)
	} else {
		cursor, err = t.Base.Exec(sql, sqldata...)
	}
	if err != nil {
		return 0, err
	}
	defer cursor.Close()

	cantidad := 0
	cursor.Next()
	err = cursor.Scan(&cantidad)
	if err != nil {
		return 0, err
	}
	return cantidad, nil
}

func (t *XTable) Min(field string, args ...interface{}) (interface{}, error) {
	// 1. analyse params
	hasconditions := false
	var conditions XConditions
	hastrx := false
	var trx *XTransaction

	for _, p := range args {
		switch p.(type) {
		case XCondition:
			hasconditions = true
			conditions = XConditions{p.(XCondition)}
		case *XCondition:
			if p != nil && p.(*XCondition) != nil {
				hasconditions = true
				conditions = XConditions{*p.(*XCondition)}
			}
		case XConditions:
			hasconditions = true
			conditions = p.(XConditions)
		case *XConditions:
			if p != nil && p.(*XConditions) != nil {
				hasconditions = true
				conditions = *p.(*XConditions)
			}
		case *XTransaction:
			trx, hastrx = p.(*XTransaction)
			if trx == nil {
				hastrx = false
			}
		}
	}

	itemdata := 1
	sqldata := make([]interface{}, 0)

	sql := "select min("
	sql += t.Prepend + field
	sql += ") from " + t.Name

	if hasconditions {
		scond, vars := conditions.CreateConditions(t, t.Base.DBType, itemdata)
		sql += " where " + scond
		for _, v := range vars {
			sqldata = append(sqldata, v)
		}
	}

	if DEBUG {
		fmt.Println(sql)
	}

	var cursor *libsql.Rows
	var err error
	if hastrx {
		cursor, err = trx.Exec(sql, sqldata...)
	} else {
		cursor, err = t.Base.Exec(sql, sqldata...)
	}
	if err != nil {
		return 0, err
	}
	defer cursor.Close()

	var min interface{}
	cursor.Next()
	err = cursor.Scan(&min)
	if err != nil {
		return 0, err
	}
	return min, nil
}

func (t *XTable) Max(field string, args ...interface{}) (interface{}, error) {
	// 1. analyse params
	hasconditions := false
	var conditions XConditions
	hastrx := false
	var trx *XTransaction

	for _, p := range args {
		switch p.(type) {
		case XCondition:
			hasconditions = true
			conditions = XConditions{p.(XCondition)}
		case *XCondition:
			if p != nil && p.(*XCondition) != nil {
				hasconditions = true
				conditions = XConditions{*p.(*XCondition)}
			}
		case XConditions:
			hasconditions = true
			conditions = p.(XConditions)
		case *XConditions:
			if p != nil && p.(*XConditions) != nil {
				hasconditions = true
				conditions = *p.(*XConditions)
			}
		case *XTransaction:
			trx, hastrx = p.(*XTransaction)
			if trx == nil {
				hastrx = false
			}
		}
	}

	itemdata := 1
	sqldata := make([]interface{}, 0)

	sql := "select max("
	sql += t.Prepend + field
	sql += ") from " + t.Name

	if hasconditions {
		scond, vars := conditions.CreateConditions(t, t.Base.DBType, itemdata)
		sql += " where " + scond
		for _, v := range vars {
			sqldata = append(sqldata, v)
		}
	}

	if DEBUG {
		fmt.Println(sql)
	}

	var cursor *libsql.Rows
	var err error
	if hastrx {
		cursor, err = trx.Exec(sql, sqldata...)
	} else {
		cursor, err = t.Base.Exec(sql, sqldata...)
	}
	if err != nil {
		return 0, err
	}
	defer cursor.Close()

	var max interface{}
	cursor.Next()
	err = cursor.Scan(&max)
	if err != nil {
		return 0, err
	}
	return max, nil
}

func (t *XTable) Avg(field string, args ...interface{}) (interface{}, error) {
	// 1. analyse params
	hasconditions := false
	var conditions XConditions
	hastrx := false
	var trx *XTransaction

	for _, p := range args {
		switch p.(type) {
		case XCondition:
			hasconditions = true
			conditions = XConditions{p.(XCondition)}
		case *XCondition:
			if p != nil && p.(*XCondition) != nil {
				hasconditions = true
				conditions = XConditions{*p.(*XCondition)}
			}
		case XConditions:
			hasconditions = true
			conditions = p.(XConditions)
		case *XConditions:
			if p != nil && p.(*XConditions) != nil {
				hasconditions = true
				conditions = *p.(*XConditions)
			}
		case *XTransaction:
			trx, hastrx = p.(*XTransaction)
			if trx == nil {
				hastrx = false
			}
		}
	}

	itemdata := 1
	sqldata := make([]interface{}, 0)

	sql := "select avg("
	sql += t.Prepend + field
	sql += ") from " + t.Name

	if hasconditions {
		scond, vars := conditions.CreateConditions(t, t.Base.DBType, itemdata)
		sql += " where " + scond
		for _, v := range vars {
			sqldata = append(sqldata, v)
		}
	}

	if DEBUG {
		fmt.Println(sql)
	}

	var cursor *libsql.Rows
	var err error
	if hastrx {
		cursor, err = trx.Exec(sql, sqldata...)
	} else {
		cursor, err = t.Base.Exec(sql, sqldata...)
	}
	if err != nil {
		return 0, err
	}
	defer cursor.Close()

	var avg interface{}
	cursor.Next()
	err = cursor.Scan(&avg)
	if err != nil {
		return 0, err
	}
	switch avg.(type) {
	case []byte:
		v, err := strconv.ParseFloat(string(avg.([]byte)), 64)
		return v, err
	case string:
		v, err := strconv.ParseFloat(avg.(string), 64)
		return v, err
	}
	return avg, nil
}

func (t *XTable) GetPrimaryKey() XFieldDef {
	for _, f := range t.Fields {
		if IsPrimaryKey(f) {
			return f
		}
	}
	return nil
}

func (t *XTable) GetField(name string) XFieldDef {
	for _, f := range t.Fields {
		if f.GetName() == name {
			return f
		}
	}
	return nil
}

func getQueryString(DB string, item int) string {
	q := ""
	switch DB {
	case DB_Postgres:
		q = fmt.Sprintf("$%d", item)
	case DB_MySQL:
		q = "?"
	}
	return q
}
