package xdominion

import (
	"database/sql"
)

type XCursor struct {
	Base        *XBase
	Transaction *XTransaction
	Query       string
	Cursor      *sql.Rows
	Columns     []*sql.ColumnType
	inuse       bool
}

func (c *XCursor) Exec(query string, args ...interface{}) (err error) {

	if c.inuse {
		c.Close()
	}

	c.Query = query
	if c.Transaction != nil {
		c.Cursor, err = c.Transaction.Exec(query, args...)
	} else if c.Base != nil {
		c.Cursor, err = c.Base.Exec(query, args...)
	}
	if err != nil {
		return err
	}

	// Get the column names and types
	c.Columns, err = c.Cursor.ColumnTypes()
	c.inuse = true
	return err
}

func (c *XCursor) Next() bool {
	return c.Cursor.Next()
}

func (c *XCursor) Read() (*XRecord, error) {

	if !c.inuse {
		return nil, nil
	}

	rec := NewXRecord()

	// Create a slice to hold the values of the current row
	values := make([]interface{}, len(c.Columns))
	for i := range values {
		values[i] = new(interface{})
	}

	// Scan the values of the current row into the slice
	if err := c.Cursor.Scan(values...); err != nil {
		return nil, err
	}

	// Iterate over the values and print the results
	for i, col := range c.Columns {
		name := col.Name()
		value := *(values[i].(*interface{}))
		rec.Set(name, value)
	}
	return rec, nil
}

func (c *XCursor) Close() {
	c.Cursor.Close()
	c.inuse = false
}
