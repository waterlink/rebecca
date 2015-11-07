package fakedriver

import (
	"fmt"

	"github.com/waterlink/rebecca/context"
	"github.com/waterlink/rebecca/field"
)

// Driver represents fake driver for tests
type Driver struct {
	whereRegistry map[string]func([]field.Field) (bool, error)
	records       map[string][][]field.Field
	maxID         int
}

// NewDriver is for creating new fake driver
func NewDriver() *Driver {
	return &Driver{
		whereRegistry: map[string]func([]field.Field) (bool, error){},
		records:       map[string][][]field.Field{},
	}
}

// Get is for fetching single record by its ID
func (d *Driver) Get(tablename string, fields []field.Field, ID field.Field) ([]field.Field, error) {
	for _, record := range d.getTable(tablename) {
		if hasField(record, ID) {
			return record, nil
		}
	}

	return nil, fmt.Errorf("Unable to find record with ID=%+v", ID.Value)
}

// Create is for creating new record. Mutates passed ID
func (d *Driver) Create(tablename string, fields []field.Field, ID *field.Field) error {
	d.maxID++
	ID.Value = d.maxID
	changeID(fields, *ID)
	d.insertTo(tablename, fields)
	return nil
}

// Update is for updating existing record
func (d *Driver) Update(tablename string, fields []field.Field, ID field.Field) error {
	records := d.getTable(tablename)
	for i, record := range records {
		if hasField(record, ID) {
			records[i] = fields
			return nil
		}
	}

	return fmt.Errorf("Unable to find record with ID=%+v", ID.Value)
}

// All is for fetching all records
func (d *Driver) All(tablename string, fields []field.Field, ctx *context.Context) ([][]field.Field, error) {
	return nil, nil
}

// Where is for fetching specific records
func (d *Driver) Where(tablename string, fields []field.Field, ctx *context.Context, where string) ([][]field.Field, error) {
	return nil, nil
}

// First is for fetching first specific record
func (d *Driver) First(tablename string, fields []field.Field, ctx *context.Context, where string) ([]field.Field, error) {
	return nil, nil
}

// Remove is for removing record by providen ID from database
func (d *Driver) Remove(tablename string, ID field.Field) error {
	return nil
}

// RegisterWhere is for registering fake where query
func (d *Driver) RegisterWhere(where string, fn func(record []field.Field) (bool, error)) {
	d.whereRegistry[where] = fn
}

func (d *Driver) ensureTable(name string) {
	_, ok := d.records[name]
	if !ok {
		d.records[name] = [][]field.Field{}
	}
}

func (d *Driver) getTable(name string) [][]field.Field {
	d.ensureTable(name)
	return d.records[name]
}

func (d *Driver) insertTo(table string, fields []field.Field) {
	d.ensureTable(table)
	d.records[table] = append(d.records[table], fields)
}

func hasField(record []field.Field, x field.Field) bool {
	for _, f := range record {
		if x == f {
			return true
		}
	}

	return false
}

func changeID(record []field.Field, ID field.Field) {
	for i, f := range record {
		if f.Primary {
			record[i] = ID
			return
		}
	}
}
