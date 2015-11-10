package rebecca

// This file contains thin exported methods related to Context only.
//
// For unexported functions see: helpers.go

import (
	"fmt"

	"github.com/waterlink/rebecca/context"
)

// Context is for storing query context
type Context struct {
	// Defines ordering of the query
	Order string

	// Defines grouping criteria of the query
	Group string

	// Defines maximum amount of records requested for the query
	Limit int

	// Defines starting record for the query
	Skip   int
	Offset int // alias of Skip

	tx interface{}
}

// GetOrder is for fetching context's Order. Used by drivers
func (c *Context) GetOrder() string {
	return c.Order
}

// GetGroup is for fetching context's Group. Used by drivers
func (c *Context) GetGroup() string {
	return c.Group
}

// GetLimit is for fetching context's Limit. Used by drivers
func (c *Context) GetLimit() int {
	return c.Limit
}

// GetSkip is for fetching context's Skip. Also it fetches Offset if present,
// hence the alias. Used by drivers
func (c *Context) GetSkip() int {
	if c.Offset > 0 {
		return c.Offset
	}
	return c.Skip
}

// GetTx is for fetching context's driver transaction state. Used by drivers
func (c *Context) GetTx() interface{} {
	return c.tx
}

// SetOrder is for setting context's Order, it creates new Context. Used by drivers
func (c *Context) SetOrder(order string) context.Context {
	ctx := c.makeCopy()
	ctx.Order = order
	return &ctx
}

// SetGroup is for setting context's Group. Used by drivers
func (c *Context) SetGroup(group string) context.Context {
	ctx := c.makeCopy()
	ctx.Group = group
	return &ctx
}

// SetLimit is for setting context's Limit. Used by drivers
func (c *Context) SetLimit(limit int) context.Context {
	ctx := c.makeCopy()
	ctx.Limit = limit
	return &ctx
}

// SetSkip is for setting context's Skip. Used by drivers
func (c *Context) SetSkip(skip int) context.Context {
	ctx := c.makeCopy()
	ctx.Skip = skip
	ctx.Offset = skip
	return &ctx
}

// All is for fetching all records
func (c *Context) All(records interface{}) error {
	meta, err := getMetadata(records)
	if err != nil {
		return err
	}

	fieldss, err := driver.All(meta.tablename, meta.fields, c)
	if err != nil {
		return fmt.Errorf("Unable to fetch all records - %s", err)
	}

	if err := populateRecordsFromFieldss(records, fieldss); err != nil {
		return fmt.Errorf("Unable to fetch all records - %s", err)
	}

	return nil
}

// Where is for fetching specific records
func (c *Context) Where(records interface{}, query string, args ...interface{}) error {
	meta, err := getMetadata(records)
	if err != nil {
		return err
	}

	fieldss, err := driver.Where(meta.tablename, meta.fields, c, query, args...)
	if err != nil {
		return fmt.Errorf("Unable to fetch specific records - %s", err)
	}

	if err := populateRecordsFromFieldss(records, fieldss); err != nil {
		return fmt.Errorf("Unable to fetch specific records - %s", err)
	}

	return nil
}

// First is for fetching only one specific record
func (c *Context) First(record interface{}, query string, args ...interface{}) error {
	meta, err := getMetadata(record)
	if err != nil {
		return err
	}

	fields, err := driver.First(meta.tablename, meta.fields, c, query, args...)
	if err != nil {
		return fmt.Errorf("Unable to fetch specific records - %s", err)
	}

	if err := setFields(record, fields); err != nil {
		return fmt.Errorf("Unable to assign fields for the record - %s", err)
	}

	return nil
}

func (c Context) makeCopy() Context {
	return c
}
