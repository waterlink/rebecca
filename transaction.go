package rebecca

// This file contains thin exported functions and methods related to
// Transaction only.
//
// For unexported functions see: helpers.go
//
// For Context see: context.go

import (
	"errors"
	"fmt"

	"github.com/waterlink/rebecca/driver"
)

// Transaction is for managing transactions for drivers that allow it
type Transaction struct {
	tx       interface{}
	finished bool
}

// Begin is for creating proper transaction
func Begin() (*Transaction, error) {
	d, lock := driver.Get()
	defer lock.Unlock()

	tx, err := d.Begin()
	if err != nil {
		return nil, fmt.Errorf("Unable to begin transaction - %s", err)
	}
	return &Transaction{
		tx: tx,
	}, nil
}

// Rollback is for rolling back the transaction
func (tx *Transaction) Rollback() {
	if tx.finished {
		return
	}

	d, lock := driver.Get()
	defer lock.Unlock()

	d.Rollback(tx.tx)
	tx.finished = true
}

// Commit is for committing the transaction
func (tx *Transaction) Commit() error {
	if tx.finished {
		return errors.New("Unable to commit transaction - Current transaction is already finished")
	}

	d, lock := driver.Get()
	defer lock.Unlock()

	if err := d.Commit(tx.tx); err != nil {
		return fmt.Errorf("Unable to commit transaction - %s", err)
	}

	tx.finished = true
	return nil
}

// Get is for fetching one record
func (tx *Transaction) Get(record interface{}, ID interface{}) error {
	return get(tx.tx, record, ID)
}

// Save is for saving one record (either creating or updating)
func (tx *Transaction) Save(record interface{}) error {
	return save(tx.tx, record)
}

// All is for fetching all records
func (tx *Transaction) All(records interface{}) error {
	ctx := tx.Context(&Context{})
	return ctx.All(records)
}

// Where is for fetching specific records
func (tx *Transaction) Where(records interface{}, where string, args ...interface{}) error {
	ctx := tx.Context(&Context{})
	return ctx.Where(records, where, args...)
}

// First is for fetching only one specific record
func (tx *Transaction) First(record interface{}, where string, args ...interface{}) error {
	ctx := tx.Context(&Context{})
	return ctx.First(record, where, args...)
}

// Remove is for removing the record
func (tx *Transaction) Remove(record interface{}) error {
	return remove(tx.tx, record)
}

// Context is for instantiating proper context for transaction
func (tx *Transaction) Context(ctx *Context) *Context {
	return &Context{
		Order: ctx.Order,
		Group: ctx.Group,
		Limit: ctx.Limit,
		Skip:  ctx.Skip,
		tx:    tx.tx,
	}
}
