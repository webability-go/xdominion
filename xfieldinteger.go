package xdominion

import (
  "fmt"
)

var fieldintegertypes = map[string]string{
    DB_Postgres: "integer",
    DB_MySQL: "integer",
/*
    DB_Base::MSSQL => array(
      DB_Field::INTEGER => "int"
    ),
    DB_Base::ORACLE => array(
      DB_Field::INTEGER => "number(16)"
    )
*/
  }

type XFieldInteger struct
{
  Name string
  Constraints XConstraints
}

// creates the name of the field with its type (to create the table)
func (f XFieldInteger)CreateField(prepend string, DB string, ifText *bool) string {
  field := prepend + f.Name
  if DB == DB_Postgres && f.IsAutoIncrement() {
    field += " serial"
  } else {
    field += " " + fieldintegertypes[DB]
  }

  /*
    if ($this->checks != null)
      $line .= $this->checks->createCheck($id.$this->name, $DB);
    else
    {
      if ($DB == DB_Base::MSSQL)
      {
        $line .= " NULL";
      }
    }
 */

  extra := ""
  if f.Constraints != nil {
    extra = f.Constraints.CreateConstraints(prepend, f.Name, DB)
  }
  return field + extra
}

// creates a string representation of the value of the field for insert/update
func (f XFieldInteger)CreateValue(v interface{}, table string, DB string, id string) string {
  return fmt.Sprint(v)
}

// creates the sequence used by the field (only autoincrement fields)
func (f XFieldInteger)CreateSequence(table string) string {
  return ""
}

// creates the index used by the field (normal, unique, multi, multi unique)
func (f XFieldInteger)CreateIndex(table string, id string, DB XBase) string {
  return ""
}

// gets the name of the field
func (f XFieldInteger)GetName() string {
  return f.Name
}

// gets the type of the field
func (f XFieldInteger)GetType() int {
  return XField_Int
}

// gets the checks of the field
func (f XFieldInteger)GetConstraints() XConstraints {
  return f.Constraints
}

// Is it autoincrement
func (f XFieldInteger)IsAutoIncrement() bool {
  if (f.Constraints != nil) {
    c := f.Constraints.Get(AI)
    if c != nil { return true }
  }
  return false;
}

