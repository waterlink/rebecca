// Package rebecca is lightweight convenience library for work with database
//
// See github README for instructions: https://github.com/waterlink/rebecca#rebecca
//
// See examples: https://godoc.org/github.com/waterlink/rebecca#pkg-examples
//
// Simple example:
//
//    type Person struct {
//            rebecca.ModelMetadata `tablename:"people"`
//
//            ID   int    `rebecca:"id" rebecca_primary:"true"`
//            Name string `rebecca:"name"`
//            Age  int    `rebecca:"age"`
//    }
//
//    // Create new record
//    p := &Person{Name: "John", Age: 34}
//    if err := rebecca.Save(p); err != nil {
//            // handle error here
//    }
//    fmt.Print(p)
//
//    // Update existing record
//    p.Name = "John Smith"
//    if err := rebecca.Save(p); err != nil {
//            // handle error here
//    }
//    fmt.Print(p)
//
//    // Get record by its primary key
//    p = &Person{}
//    if err := rebecca.Get(25, p); err != nil {
//            // handle error here
//    }
//    fmt.Print(p)
package rebecca

import "github.com/waterlink/rebecca/driver"

// This file contains thin exported functions only.
//
// For unexported functions see: helpers.go
//
// For Context see: context.go

// SetupDriver is for configuring database driver
func SetupDriver(d driver.Driver) {
	driver.SetupDriver(d)
}

// Get is for fetching one record
func Get(ID interface{}, record interface{}) error {
	return get(nil, ID, record)
}

// All is for fetching all records
func All(records interface{}) error {
	ctx := &Context{}
	return ctx.All(records)
}

// Where is for fetching specific records
func Where(records interface{}, where string, args ...interface{}) error {
	ctx := &Context{}
	return ctx.Where(records, where, args...)
}

// First is for fetching only one specific record
func First(record interface{}, where string, args ...interface{}) error {
	ctx := &Context{}
	return ctx.First(record, where, args...)
}

// Save is for saving one record (either creating or updating)
func Save(record interface{}) error {
	return save(nil, record)
}

// Remove is for removing the record
func Remove(record interface{}) error {
	return remove(nil, record)
}
