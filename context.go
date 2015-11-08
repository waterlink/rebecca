package rebecca

import (
	"fmt"

	"github.com/waterlink/rebecca/context"
)

// Context is for storing query context
type Context struct {
	Order string
	Group string
	Limit int
	Skip  int

	tx interface{}
}

// GetOrder is for fetching context's Order
func (c *Context) GetOrder() string {
	return c.Order
}

// GetGroup is for fetching context's Group
func (c *Context) GetGroup() string {
	return c.Group
}

// GetLimit is for fetching context's Limit
func (c *Context) GetLimit() int {
	return c.Limit
}

// GetSkip is for fetching context's Skip
func (c *Context) GetSkip() int {
	return c.Skip
}

// GetTx is for fetching context's driver transaction state
func (c *Context) GetTx() interface{} {
	return c.tx
}

// SetOrder is for setting context's Order, it creates new Context
func (c *Context) SetOrder(order string) context.Context {
	return &Context{
		Order: order,
		Group: c.Group,
		Limit: c.Limit,
		Skip:  c.Skip,
		tx:    c.tx,
	}
}

// SetGroup is for setting context's Group
func (c *Context) SetGroup(group string) context.Context {
	return &Context{
		Order: c.Order,
		Group: group,
		Limit: c.Limit,
		Skip:  c.Skip,
		tx:    c.tx,
	}
}

// SetLimit is for setting context's Limit
func (c *Context) SetLimit(limit int) context.Context {
	return &Context{
		Order: c.Order,
		Group: c.Group,
		Limit: limit,
		Skip:  c.Skip,
		tx:    c.tx,
	}
}

// SetSkip is for setting context's Skip
func (c *Context) SetSkip(skip int) context.Context {
	return &Context{
		Order: c.Order,
		Group: c.Group,
		Limit: c.Limit,
		Skip:  skip,
		tx:    c.tx,
	}
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
