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
- Conditions
- Sub Queries
- Group and report funcions
- Synchro to upgrade DB tables and fields
- MySQL, Oracle, other DBs

Version Changes Control
=======================

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

The manual is still under construction.

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


---
