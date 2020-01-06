@UTF-8

XDominion for GO v0
=============================

The XDominion library is used to build object instead of queries to access any database. Queries are build on demand based on the type of database. 
Do not write queries anymore, but use objects.

Because you need to compile, add the database drivers directly into the code, 
XDominion support only posgresql in this version. More databases will be supported later.

XDominion is the Go adaptation of PHP7-Dominion libraries: a full database abstraction layer

TO DO:
======

- text, float, time, date, lob fields
- Joins
- Sub Queries
- Group and report functions
- Synchro to upgrade DB tables and fields
- Oracle, Informix, Mongo, other DBs

Version Changes Control
=======================

v0.1.2 - 2020-01-06
-----------------------
- Corrected an error in Select, Update, Delete to use int32, int64, float32 as values
- Added functions Min, Max, Avg

v0.1.1 - 2019-12-17
-----------------------
- Corrected an error in Upsert func that was inserting a 0 even if the key was present into the record

v0.1.0 - 2019-12-06
-----------------------
- Added new XTable Language functionality to know the default language of data into a table
- gofmt code formated before pushing a change on github

v0.0.15 - 2019-10-11
-----------------------
- Added record conversions for float values, from string, float32, int, etc
- Added Upsert function in table (update or insert if not exists)

v0.0.14 - 2019-06-25
-----------------------
- Added Clone on XRecord and XRecords to meet definition of xcore.XDatasetDef and xcore.XDatasetCollectionDef

v0.0.13 - 2019-06-19
-----------------------
- Error corrected on XTable.SelectAll. It was not working as expected
- Error corrected on XRecord.GetString. was returning "<nil>" instead of "" when the database field was null

v0.0.12 - 2019-03-06
-----------------------
- Support for Time functions added in the XRecord (instanciated from XDatasetDef)

v0.0.11 - 2019-03-01
-----------------------
- Many correction on Mysql support to make it work correctly
- Removed GetValue function from FieldDef
- "?" implemented for fields, conditions and having queries (in select, insert, update, delete statements)
- Values are directly passed to the query with "?", not a string representation of them

v0.0.10 - 2019-02-18
-----------------------
- Support for MySQL added
- Queries and conditions now uses "?" or "$x" for parameters
- Orderby implemented
- like and ilike implemented for text fields
- fields.GetValue function added when the code needs the raw string value (not like CreateValue where then value is created with ' for strings)

v0.0.9 - 2019-02-15
-----------------------
- New funcion for field: GetValue created
- Error corrected on conflict between CreateValue (with ' for strings) and GetValue (for use with $d to inject into queries)

v0.0.8 - 2019-02-14
-----------------------
- XCondition works with string queries (not yet with "?" parameters)
- Correction done on CreateValue for string fields (text, varchar, dates)

v0.0.7 - 2019-02-05
-----------------------
- Added XOrderBy, XOrder structures
- Added XGroupBy, XGroup structures
- Added XHaving structures
- Error corrected on xrecord.GetTime when the field comes NIL from database
- Added DEBUG main xdominion global variable. Set to true to print all the built queries

V0.0.6 - 2019-01-15
-----------------------
- Added conversion between types con Get* functions
- XTable.Count implemented

V0.0.5 - 2019-01-15
-----------------------
- XTable.Update implemented
- XTable.Delete implemented

V0.0.4 - 2019-01-06
-----------------------
- Modify XRecord to match XDataset last version (xcore 0.0.4)
- Modify XRecords to match XDatasetCollection last version (xcore 0.0.4)

V0.0.3 - 2018-12-19
-----------------------
- XField added: Float, Date, DateTime, Text, partially implemented
- XConditions and XCondition added, partially implemented
- XConstraints and XConstaint added, partially implemented
- XFieldSet added, partially implemented
- XOrderBy added, partially implemented

V0.0.2 - 2018-12-17
-----------------------
- Postgres implementation
- XTable created. select and insert partially done
- XField created, Integer and VarChar partially done
- XRecord created
- XRecords created
- XCursor created

V0.0.1 - 2018-11-14
-----------------------
- First commit, Eventually does not work yet
- Base object done


Manual:
=======================

1. Overview
------------------------

XDominion is a dataase abstraction layer, to build and use objects of data instead of building SQL queries. 
The code is portable between databases with changing the implementation, since you don't use direct incompatible SQL sentences.

The library is build over 3 main objects:
- XBase: database connector and cursors to build queries and manipulation language
- - Other included objects: XCursor
- XTable: the table definition, data access function/structures and definition manipulation language
- - Other included objects: XField*, XConstraints, XContraint, XOrderby, XConditions, XCondition
- XRecord: the results and data to interchange with the database
- - Other included objects: XRecords


Examples:

Some code to start working:

Creates the connector to the database and connect:

```
  base := &xdominion.XBase{
    DBType: xdominion.DB_Postgres,
    Username: "username",
    Password: "password",
    Database: "test",
    Host: xdominion.DB_Localhost,
    SSL: false,
  }
  base.Logon()
```

Executes a query:

```
  q, err := base.Exec("drop table test")
  if (err != nil) {
    fmt.Println(err)
  }
  q.Close()
```

Creates a table definition:

```
t := xdominion.NewXTable("test", "t_")
t.AddField(xdominion.XFieldInteger{Name: "f1", Constraints: xdominion.XConstraints{
                                                  xdominion.XConstraint{Type: xdominion.PK},
                                                  xdominion.XConstraint{Type: xdominion.AI},
                                               } })   // ai, pk
t.AddField(xdominion.XFieldVarChar{Name: "f2", Size: 20, Constraints: xdominion.XConstraints{
                                                  xdominion.XConstraint{Type: xdominion.NN},
                                               } })
t.AddField(xdominion.XFieldText{Name: "f3"})
t.AddField(xdominion.XFieldDate{Name: "f4"})
t.AddField(xdominion.XFieldDateTime{Name: "f5"})
t.AddField(xdominion.XFieldFloat{Name: "f6"})
t.SetBase(base)
```

Synchronize the table with DB (create it if it does not exist)

```
  err = t.Synchronize()
  if (err != nil) {
    fmt.Println(err)
  }
```

Some Insert:

```
  res1, err := tb.Insert(xdominion.XRecord{"f1": 1, "f2": "Data line 1",})
  if (err != nil) {
    fmt.Println(err)
  }
  fmt.Println(res1)  // res1 is the primary key
```

With an error (f2 is mandatory based on table definition):

```
  res21, err := tb.Insert(xdominion.XRecord{"f1": 2, "f3": "test",})
  if (err != nil) {
    fmt.Println(err)
  }
  fmt.Println(res21)
```

General query (select ALL):
```
  res3, err := tb.Select()
  if err != nil {
    fmt.Println(err)
  } else {
    for _, x := range res3.(xdominion.XRecords) {
      fmt.Println(x)
    }
  }
```

Query by Key:

```
  res4, err := tb.Select(1)
  if err != nil {
    fmt.Println(err)
  } else {
    switch res4.(type) {
      case xdominion.XRecord:
        fmt.Println(res4)
      case xdominion.XRecords:
        for _, x := range res4.(xdominion.XRecords) {
          fmt.Println(x)
        }
    }
  }
```
  
Query by Where:
  
```
  res5, err := tb.Select(xdominion.XConditions{xdominion.NewXCondition("f1", "=", 1), xdominion.NewXCondition("f2", "like", "lin", "and")})
  if err != nil {
    fmt.Println(err)
  } else {
    switch res5.(type) {
      case xdominion.XRecord:
        fmt.Println(res5)
      case xdominion.XRecords:
        for _, x := range res5.(xdominion.XRecords) {
          fmt.Println(x)
        }
    }
  }
```


2. Reference
------------------------

XBase
-----

XTable
------

XRecord
-------



---
