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


---
