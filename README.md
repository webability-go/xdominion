
XDominion for GO v0
=============================

[![Go Report Card](https://goreportcard.com/badge/github.com/webability-go/xdominion)](https://goreportcard.com/report/github.com/webability-go/xdominion)
[![GoDoc](https://godoc.org/github.com/webability-go/xdominion?status.png)](https://godoc.org/github.com/webability-go/xdominion)


The XDominion library is used to build object instead of queries to access any database. Queries are build on demand based on the type of database.
Do not write queries anymore, but use objects.

Because you need to compile, add the database drivers directly into the code,
XDominion support only posgresql and MySQL in this version. More databases will be supported later.
If you need a not yet supported database, please open a ticket on github.com.

The library provides a set of high-level APIs for interacting with databases. It allows developers to map database tables to Go structs, allowing them to interact with the database using objects. The library also provides an intuitive and chainable API for querying the database, similar to the structure of SQL statements, but without requiring developers to write SQL code directly.

xdominion uses a set of interfaces to abstract the database operations, making it easy to use different database backends with the same code. The library supports transactions, allowing developers to perform multiple database operations in a single transaction.

The xdominion library uses reflection to map Go structs to database tables, and also allows developers to specify custom column names and relationships between tables.

Overall, xdominion provides a simple and intuitive way to interact with databases using objects and abstracts the underlying database implementation. It is a well-designed library with a clear API and support for multiple database backends.

XDominion needs Go v1.17+





Version Changes Control
=======================

v0.5.2 - 2025-02-14
-----------------------
- Aggregation capability added to the XTable and related objects. (GroupBy, Group, Select <functions>, count, avg, min, max, now, etc)
- Upgraded golang.org/x/net for a security patch to v0.33.0
 
v0.5.1 - 2024-04-09
-----------------------
- Correction of a bug into manuals

v0.5.0 - 2023-05-02
-----------------------
- XCursor is now implemented and makes a simple query to the database but return results as an *XRecord.
- XCursor is complatible with XBase and XTransactions
- Official documentation mounted for https://pkg.go.dev/

v0.4.3 - 2022-11-22
-----------------------
- Added GetInt for a string with automatic conversion if possible from an XRecord

v0.4.2 - 2022-06-13
-----------------------
- Correction of a bug on DoSelect with limits, a 0 limit is not a limit

v0.4.1 - 2022-01-19
-----------------------
- XConditions, XCondition, XOrder, XOrderBy, XFieldSet can now be cloned with <object>.Clone()

v0.4.0 - 2021-12-10
-----------------------
- Now all the xtable functions (Select, Update, Delete, Count, etc) can accept pointers parameters and nil casted parameters (for instance *XCondition or *XOrderBy than can be nil too).

v0.3.3 - 2021-01-20
-----------------------
- Correction to support nil transactions (no transaction even if the parameter is passed with a nil value) into select type queries (select, min, max, avg, count, ...)

v0.3.2 - 2021-01-17
-----------------------
- Implementation of transactions into select type queries (select, min, max, avg, count, ...)

v0.3.1 - 2020-12-04
-----------------------
- Implementation of indexes creation during the table synchronization. Supports now index, unique index, multiple index and multiple unique index.

v0.3.0 - 2020-11-10
-----------------------
- Implementation of transactions, new XTransaction object and functions to create a transaction, commit or rollback it.

v0.2.3 - 2020-06-03
-----------------------
- Upgrade to xcore/v2
- Modularization with go.mod

v0.2.2 - 2020-06-03
-----------------------
- Bug Corrected on Clonation of XRecord, it now consider XRecords (via interface Clone() XDatasetCollectionDef) as possible subset to clone too.

v0.2.1 - 2020-02-11
-----------------------
- Bug Corrected on String and GoString of XRecord and XRecords

v0.2.0 - 2020-02-10
-----------------------
- Modification to XRecord and XRecords to meet xcore v1.0.0 (.String and .GoString functions added, .Stringify function removed)

v0.1.3 - 2020-01-08
-----------------------
- Corrected an error in Insert, Update to use XRecordDef and XRecordsDef instead of XRecord and XRecords to be widely compatible with any entry parameter
- All functions will auto-identify if parameters are XRecord, *XRecord, XRecords or *XRecords

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


---
