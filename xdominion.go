// Copyright Philippe Thomassigny 2004-2023.
// Use of this source code is governed by a MIT licence.
// license that can be found in the LICENSE file.
//
// XDominion for GO v0
// =============================
//
// xdominion is a Go library for creating a database layer that abstracts the underlying database implementation and allows developers to interact with the database using objects rather than SQL statements. It supports multiple database backends, including PostgreSQL, MySQL, SQLite, and Microsoft SQL Server, among others.
//
// If you need a not yet supported database, please open a ticket on github.com.
//
// The library provides a set of high-level APIs for interacting with databases. It allows developers to map database tables to Go structs, allowing them to interact with the database using objects. The library also provides an intuitive and chainable API for querying the database, similar to the structure of SQL statements, but without requiring developers to write SQL code directly.
//
// xdominion uses a set of interfaces to abstract the database operations, making it easy to use different database backends with the same code. The library supports transactions, allowing developers to perform multiple database operations in a single transaction.
//
// The xdominion library uses reflection to map Go structs to database tables, and also allows developers to specify custom column names and relationships between tables.
//
// Overall, xdominion provides a simple and intuitive way to interact with databases using objects and abstracts the underlying database implementation. It is a well-designed library with a clear API and support for multiple database backends.
//
//
// 1. Overview
// ------------------------
//
// XDominion is a database abstraction layer, to build and use objects of data instead of building SQL queries.
// The code is portable between databases with changing the implementation, since you don't use direct incompatible SQL sentences.
//
// The library is build over 3 main objects:
// - XBase: database connector and cursors to build queries and manipulation language
// - - Other included objects: XCursor
// - XTable: the table definition, data access function/structures and definition manipulation language
// - - Other included objects: XField*, XConstraints, XContraint, XOrderby, XConditions, XCondition
// - XRecord: the results and data to interchange with the database
// - - Other included objects: XRecords
//
//
// 2. Some example code to start working rapidly:
// ------------------------
//
// Creates the connector to the database and connect:
//
// ```
//   base := &xdominion.XBase{
//     DBType: xdominion.DB_Postgres,
//     Username: "username",
//     Password: "password",
//     Database: "test",
//     Host: xdominion.DB_Localhost,
//     SSL: false,
//   }
//   base.Logon()
// ```
//
// Executes a direct query:
//
// ```
//   q, err := base.Exec("drop table test")
//   if (err != nil) {
//     fmt.Println(err)
//   }
//   defer q.Close()
//
//   // With a select and an XCursor:
//
//	// 1. Create a cursor
//	cursor := base.NewXCursor()
//	err := cursor.Exec("select * from test")
//	fmt.Println(err, cursor)
//
//	for cursor.Next() {
//		data, err := cursor.Read()  // data is an *XRecord
//		fmt.Println(err, data)
//	}
//	cursor.Close()
// ```
//
// Creates a table definition:
//
// ```
// t := xdominion.NewXTable("test", "t_")
// t.AddField(xdominion.XFieldInteger{Name: "f1", Constraints: xdominion.XConstraints{
//                                                   xdominion.XConstraint{Type: xdominion.PK},
//                                                   xdominion.XConstraint{Type: xdominion.AI},
//                                                } })   // ai, pk
// t.AddField(xdominion.XFieldVarChar{Name: "f2", Size: 20, Constraints: xdominion.XConstraints{
//                                                   xdominion.XConstraint{Type: xdominion.NN},
//                                                } })
// t.AddField(xdominion.XFieldText{Name: "f3"})
// t.AddField(xdominion.XFieldDate{Name: "f4"})
// t.AddField(xdominion.XFieldDateTime{Name: "f5"})
// t.AddField(xdominion.XFieldFloat{Name: "f6"})
// t.SetBase(base)
// ```
//
// Synchronize the table with DB (create it if it does not exist)
//
// ```
//   err = t.Synchronize()
//   if (err != nil) {
//     fmt.Println(err)
//   }
// ```
//
// Some Insert:
//
// ```
//   res1, err := tb.Insert(xdominion.XRecord{"f1": 1, "f2": "Data line 1",})
//   if (err != nil) {
//     fmt.Println(err)
//   }
//   fmt.Println(res1)  // res1 is the primary key
// ```
//
// With an error (f2 is mandatory based on table definition):
//
// ```
//   res21, err := tb.Insert(xdominion.XRecord{"f1": 2, "f3": "test",})
//   if (err != nil) {
//     fmt.Println(err)
//   }
//   fmt.Println(res21)
// ```
//
// General query (select ALL):
// ```
//   res3, err := tb.Select()
//   if err != nil {
//     fmt.Println(err)
//   } else {
//     for _, x := range res3.(xdominion.XRecords) {
//       fmt.Println(x)
//     }
//   }
// ```
//
// Query by Key:
//
// ```
//   res4, err := tb.Select(1)
//   if err != nil {
//     fmt.Println(err)
//   } else {
//     switch res4.(type) {
//       case xdominion.XRecord:
//         fmt.Println(res4)
//       case xdominion.XRecords:
//         for _, x := range res4.(xdominion.XRecords) {
//           fmt.Println(x)
//         }
//     }
//   }
// ```
//
// Query by Where:
//
// ```
//   res5, err := tb.Select(xdominion.XConditions{xdominion.NewXCondition("f1", "=", 1), xdominion.NewXCondition("f2", "like", "lin", "and")})
//   if err != nil {
//     fmt.Println(err)
//   } else {
//     switch res5.(type) {
//       case xdominion.XRecord:
//         fmt.Println(res5)
//       case xdominion.XRecords:
//         for _, x := range res5.(xdominion.XRecords) {
//           fmt.Println(x)
//         }
//     }
//   }
// ```
//
// Transactions:
//
// ```
// tx, err := base.BeginTransaction()
// res1, err := tb.Insert(XRecord{"f1": 5, "f2": "Data line 1"}, tx)
// res2, err := tb.Update(2, XRecord{"f1": 5, "f2": "Data line 1"}, tx)
// res3, err := tb.Delete(3, tx)
// // Note that the transaction is always passed as a parameter to the insert, update, delete operations
// if err != nil {
//   tx.Rollback()
//   return err
// }
// tx.Commit()
// ```
//
//
// 3. Reference
// ------------------------
//
// XBase
// -----
//
// The xbase package in xdominion provides a set of functions for working with relational databases in Go. Here is a reference manual for the package:
//
// Constants
// VERSION: A constant string that represents the version of XDominion.
// DB_Postgres: A constant string that represents the PostgreSQL database.
// DB_MySQL: A constant string that represents the MySQL database.
// DB_Localhost: A constant string that represents the local host.
// Variables
// DEBUG: A boolean variable used to enable/disable debug mode.
// Structs
// XBase
// DB: A pointer to an instance of sql.DB, representing the database connection.
// Logged: A boolean indicating whether the database connection has been established.
// DBType: A string representing the type of database being used.
// Username: A string representing the username for the database connection.
// Password: A string representing the password for the database connection.
// Database: A string representing the name of the database being connected to.
// Host: A string representing the host for the database connection.
// SSL: A boolean indicating whether to use SSL for the database connection.
// Logger: A pointer to a logger for debugging purposes.
// XTransaction
// DB: A pointer to an instance of XBase, representing the database connection.
// TX: A pointer to an instance of sql.Tx, representing a transaction.
// Functions
// Logon()
// The Logon() function establishes a connection to the database.
//
// go
// Copy code
// func (b *XBase) Logon()
// Logoff()
// The Logoff() function closes the database connection.
//
// go
// Copy code
// func (b *XBase) Logoff()
// Exec()
// The Exec() function executes a SQL query on the database and returns a cursor.
//
// go
// Copy code
// func (b *XBase) Exec(query string, args ...interface{}) (*sql.Rows, error)
// Cursor()
// The Cursor() function returns a new instance of Cursor, which provides methods for working with database records.
//
//
//
//
//
// go
// Copy code
// package main
//
// import (
// 	"fmt"
//
// 	"github.com/webability-go/xdominion"
// )
//
// func main() {
// 	// Create a new database connection
// 	base := &xdominion.XBase{
// 		DBType:   xdominion.DB_Postgres,
// 		Username: "username",
// 		Password: "password",
// 		Database: "test",
// 		Host:     xdominion.DB_Localhost,
// 		SSL:      false,
// 	}
//
// 	// Connect to the database
// 	base.Logon()
//
// 	// Execute a query
// 	query := "INSERT INTO users (name, email) VALUES ($1, $2)"
// 	_, err := base.Exec(query, "John Doe", "john.doe@example.com")
// 	if err != nil {
// 		fmt.Println("Error executing query:", err)
// 	}
//
// 	// Close the database connection
// 	base.Logoff()
// }
// In this example, we first create a new instance of the xdominion.XBase struct with the connection details to the database we want to connect to. We then call the Logon() method of the XBase struct to establish a connection to the database.
//
// Next, we define an SQL query to insert a new user into the users table, and then call the Exec() method of the XBase struct with the query and the values we want to insert. The Exec() function returns a cursor, which we don't need in this example, so we ignore it using the blank identifier (_).
//
// If there's an error executing the query, we print an error message to the console. Finally, we close the database connection by calling the Logoff() method of the XBase struct.
//
// Note that this is just a simple example, and you should always make sure to properly handle errors and sanitize user input when working with databases.
//
//
// package main
//
// import (
// 	"fmt"
//
// 	"github.com/webability-go/xdominion"
// )
//
// func main() {
// 	// Create a new database connection
// 	base := &xdominion.XBase{
// 		DBType:   xdominion.DB_Postgres,
// 		Username: "username",
// 		Password: "password",
// 		Database: "test",
// 		Host:     xdominion.DB_Localhost,
// 		SSL:      false,
// 	}
//
// 	// Connect to the database
// 	base.Logon()
//
// 	// Execute a query
// 	query := "SELECT name, email FROM users WHERE id=$1"
// 	rows, err := base.Exec(query, 1)
// 	if err != nil {
// 		fmt.Println("Error executing query:", err)
// 	}
//
// 	// Iterate over the rows and print the results
// 	for rows.Next() {
// 		var name, email string
// 		err := rows.Scan(&name, &email)
// 		if err != nil {
// 			fmt.Println("Error reading row:", err)
// 			continue
// 		}
// 		fmt.Println("Name:", name)
// 		fmt.Println("Email:", email)
// 	}
//
// 	// Close the rows and database connection
// 	rows.Close()
// 	base.Logoff()
// }
// In this example, we first create a new instance of the xdominion.XBase struct with the connection details to the database we want to connect to. We then call the Logon() method of the XBase struct to establish a connection to the database.
//
// Next, we define an SQL query to select a user from the users table with the id equal to 1. We then call the Exec() method of the XBase struct with the query and the value we want to use for the id parameter. The Exec() function returns a cursor that we can iterate over to get the results of the query.
//
// We use a for loop to iterate over the rows returned by the Exec() function. Inside the loop, we use the Scan() method of the rows object to read the values of the name and email columns into variables. We then print the values of these variables to the console.
//
// If there's an error executing the query or reading a row, we print an error message to the console. Finally, we close the rows object and the database connection by calling the Close() and Logoff() methods of the XBase struct, respectively.
//
// Note that this is just a simple example, and you should always make sure to properly handle errors and sanitize user input when working with databases.
//
//
//
//
//
//
//
//
//
// go
// Copy code
// func (b *XBase) Cursor() *Cursor
// BeginTransaction()
// The BeginTransaction() function starts a new transaction on the database.
//
// go
// Copy code
// func (b *XBase) BeginTransaction() (*XTransaction, error)
// Commit()
// The Commit() function commits a transaction to the database.
//
// go
// Copy code
// func (t *XTransaction) Commit() error
// Rollback()
// The Rollback() function rolls back a transaction on the database.
//
// go
// Copy code
// func (t *XTransaction) Rollback() error
// Notes
// The Logon() function must be called before using any other functions in the xbase package.
// The Logoff() function should be called when finished using the database connection.
// The Exec() function should be used for executing arbitrary SQL queries.
// The Cursor() function should be used for performing CRUD operations on database records.
// The BeginTransaction(), Commit(), and Rollback() functions should be used for transactions.
// Note that this is just a brief overview of the xbase package. For more information and examples, please refer to the documentation in the xdominion GitHub repository: https://github.com/webability-go/xdominion.
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
// Create a new instance of the xdominion.XBase struct, which represents a database connection. The XBase struct provides methods for interacting with the database, such as querying, inserting, updating, and deleting records.
//
// base := &xdominion.XBase{
//     DBType: xdominion.DB_Postgres,
//     Username: "username",
//     Password: "password",
//     Database: "test",
//     Host: xdominion.DB_Localhost,
//     SSL: false,
// }
//
// In this example, &xdominion.XBase{} is the instance of the XBase struct, and the properties of the struct are set to the database connection details. The DBType property specifies the type of database being used, Username and Password specify the username and password for the database connection, Database specifies the name of the database being connected to, Host specifies the host for the database connection, and SSL specifies whether to use SSL for the database connection.
//
// Use the Logon() method of the XBase struct to connect to the database.
//
// base.Logon()
//
// The Logon() method establishes a connection to the database using the details provided in the XBase struct.
//
// Note that this is just a simple example, and the XBase library provides many more features for working with databases using objects. You can find more information and examples in the xdominion GitHub repository: https://github.com/webability-go/xdominion.
//
//
// XTable definition
// -----------------
//
// XTable operations
// -----------------
//
// XRecord
// -------
//
// XRecords
// --------
//
// Conditions
// ----------
//
// Orderby
// -------
//
// Fields
// ------
//
// Limits
// ------
//
// Groupby
// -------
//
// Having
// ------
//
//
// */
package xdominion

// VERSION is the used version number of the XDominion library.
const VERSION = "0.5.0"

// DEBUG is the flag to activate debugging on the library.
// if DEBUG is set to TRUE, DEBUG indicates to the XDominion libraries to log a trace of queries and functions called, with most important parameters.
// DEBUG can be set to true or false dynamically to trace only parts of code on demand.
var DEBUG bool = false
