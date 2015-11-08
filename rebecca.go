// Package rebecca is lightweight convenience library for work with database
package rebecca

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/waterlink/rebecca/context"
	"github.com/waterlink/rebecca/field"
)

var (
	driver Driver
)

// Driver is for abstracting interaction with specific database
type Driver interface {
	Get(tx interface{}, tablename string, fields []field.Field, ID field.Field) ([]field.Field, error)
	Create(tx interface{}, tablename string, fields []field.Field, ID *field.Field) error
	Update(tx interface{}, tablename string, fields []field.Field, ID field.Field) error
	All(tablename string, fields []field.Field, ctx context.Context) ([][]field.Field, error)
	Where(tablename string, fields []field.Field, ctx context.Context, where string, args ...interface{}) ([][]field.Field, error)
	First(tablename string, fields []field.Field, ctx context.Context, where string, args ...interface{}) ([]field.Field, error)
	Remove(tx interface{}, tablename string, ID field.Field) error
	HasTransactions() bool
	Begin() (interface{}, error)
	Rollback(tx interface{})
	Commit(tx interface{}) error
}

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

// Transaction is for managing transactions for drivers that allow it
type Transaction struct {
	tx       interface{}
	finished bool
}

// Begin is for creating proper transaction
func Begin() (*Transaction, error) {
	tx, err := driver.Begin()
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

	driver.Rollback(tx.tx)
	tx.finished = true
}

// Commit is for committing the transaction
func (tx *Transaction) Commit() error {
	if tx.finished {
		return errors.New("Unable to commit transaction - Current transaction is already finished")
	}

	if err := driver.Commit(tx.tx); err != nil {
		return fmt.Errorf("Unable to commit transaction - %s", err)
	}

	tx.finished = true
	return nil
}

// Get is for fetching one record
func (tx *Transaction) Get(ID interface{}, record interface{}) error {
	return get(tx.tx, ID, record)
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

// ModelMetadata is for storing any metadata for the whole model
type ModelMetadata struct{}

type metadata struct {
	tablename string
	fields    []field.Field
	primary   field.Field
}

// Get is for fetching one record
func Get(ID interface{}, record interface{}) error {
	return get(nil, ID, record)
}

// Save is for saving one record (either creating or updating)
func Save(record interface{}) error {
	return save(nil, record)
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

// Remove is for removing the record
func Remove(record interface{}) error {
	return remove(nil, record)
}

// SetupDriver is for setting up driver manually
func SetupDriver(d Driver) {
	driver = d
}

func getMetadata(record interface{}) (metadata, error) {
	meta, err := fetchMetadata(record)
	if err != nil {
		return metadata{}, fmt.Errorf(
			"Unable to fetch record's metadata - All records are required to embed rebecca.ModelMetadata - %s",
			err,
		)
	}
	return meta, nil
}

func typeHasElem(ty reflect.Type) bool {
	return ty.Kind() == reflect.Ptr ||
		ty.Kind() == reflect.Interface ||
		ty.Kind() == reflect.Slice
}

func valueHasElem(v reflect.Value) bool {
	return v.Kind() == reflect.Ptr ||
		v.Kind() == reflect.Interface
}

func fetchMetadata(record interface{}) (metadata, error) {
	missingMetadata := metadata{}

	ty := reflect.TypeOf(record)
	for typeHasElem(ty) {
		ty = ty.Elem()
	}

	if ty.Kind() != reflect.Struct {
		return missingMetadata, fmt.Errorf("Rebecca's model is required to be struct, but got: %+v", record)
	}

	metaField, ok := ty.FieldByName("ModelMetadata")
	if !ok {
		return missingMetadata, fmt.Errorf("Rebecca's model is required to embed rebecca.ModelMetadata")
	}

	metaType := metaField.Type
	if metaType.PkgPath() != "github.com/waterlink/rebecca" || metaType.Name() != "ModelMetadata" {
		return missingMetadata, fmt.Errorf("Rebecca's model is required to embed rebecca.ModelMetadata")
	}

	meta := metadata{}

	metaTag := metaField.Tag
	tablename := metaTag.Get("tablename")
	if tablename == "" {
		return missingMetadata, fmt.Errorf("tablename tag metadata is missing on ModelMetadata")
	}

	meta.tablename = tablename

	fieldCount := ty.NumField()
	for i := 0; i < fieldCount; i++ {
		f := ty.Field(i)
		if f.Name == "ModelMetadata" {
			continue
		}

		metaField := field.Field{
			Name:       f.Name,
			Ty:         f.Type,
			DriverName: driverName(f),
			Primary:    isPrimary(f),
		}

		if metaField.Primary {
			meta.primary = metaField
		}

		meta.fields = append(meta.fields, metaField)
	}

	return meta, nil
}

func driverName(field reflect.StructField) string {
	name := field.Tag.Get("rebecca")
	if name == "" {
		name = field.Name
	}
	return name
}

func isPrimary(field reflect.StructField) bool {
	return field.Tag.Get("rebecca_primary") == "true"
}

func setFields(record interface{}, fields []field.Field) error {
	for _, f := range fields {
		if err := assignField(record, f); err != nil {
			return err
		}
	}

	return nil
}

func assignField(record interface{}, f field.Field) error {
	v := reflect.ValueOf(record)
	for valueHasElem(v) {
		v = v.Elem()
	}

	if _, ok := v.Type().FieldByName(f.Name); !ok {
		return fmt.Errorf("Field %s not found on record %+v", f.Name, record)
	}

	vf := v.FieldByName(f.Name)
	if !vf.CanSet() {
		return fmt.Errorf(
			"Unable to set field %s on record %+v. It is required to be exported and addressable",
			f.Name,
			record,
		)
	}

	vf.Set(reflect.ValueOf(f.Value))

	return nil
}

func fieldsFor(meta *metadata, record interface{}) ([]field.Field, error) {
	v := reflect.ValueOf(record)
	for valueHasElem(v) {
		v = v.Elem()
	}

	ty := v.Type()

	fields := []field.Field{}
	for _, f := range meta.fields {
		if _, ok := ty.FieldByName(f.Name); !ok {
			return nil, fmt.Errorf("Field %s not found on record %+v", f.Name, record)
		}

		itsField := f
		itsField.Value = v.FieldByName(f.Name).Interface()
		fields = append(fields, itsField)
	}

	return fields, nil
}

func populateFieldValue(record interface{}, f *field.Field) error {
	v := reflect.ValueOf(record)
	for valueHasElem(v) {
		v = v.Elem()
	}

	ty := v.Type()

	if _, ok := ty.FieldByName(f.Name); !ok {
		return fmt.Errorf("Field %s not found on record %+v", f.Name, record)
	}

	f.Value = v.FieldByName(f.Name).Interface()
	return nil
}

func isNewRecord(record interface{}, ID field.Field) (bool, error) {
	v := reflect.ValueOf(record)
	for valueHasElem(v) {
		v = v.Elem()
	}

	ty := v.Type()

	if _, ok := ty.FieldByName(ID.Name); !ok {
		return false, fmt.Errorf("Field %s not found on record %+v", ID.Name, record)
	}

	f := v.FieldByName(ID.Name)
	return f.Interface() == reflect.Zero(f.Type()).Interface(), nil
}

func zeroValueOf(value interface{}) interface{} {
	ty := reflect.TypeOf(value)
	for typeHasElem(ty) {
		ty = ty.Elem()
	}

	return reflect.New(ty).Interface()
}

func populateRecordsFromFieldss(records interface{}, fieldss [][]field.Field) error {
	for _, fields := range fieldss {
		record := zeroValueOf(records)
		if err := setFields(&record, fields); err != nil {
			return fmt.Errorf("Unable to assign fields for new record - %s", err)
		}
		v := reflect.ValueOf(records).Elem()
		v.Set(reflect.Append(v, reflect.ValueOf(record).Elem()))
	}

	return nil
}

func get(tx interface{}, ID interface{}, record interface{}) error {
	meta, err := getMetadata(record)
	if err != nil {
		return err
	}

	idField := meta.primary
	idField.Value = ID

	fields, err := driver.Get(tx, meta.tablename, meta.fields, idField)
	if err != nil {
		return fmt.Errorf("Unable to find record - %s", err)
	}

	if err := setFields(record, fields); err != nil {
		return fmt.Errorf("Unable to construct found record - %s", err)
	}

	return nil
}

func save(tx interface{}, record interface{}) error {
	meta, err := getMetadata(record)
	if err != nil {
		return err
	}

	fields, err := fieldsFor(&meta, record)
	if err != nil {
		return fmt.Errorf("Unable to fetch fields for record %+v", record)
	}

	idField := meta.primary
	isNew, err := isNewRecord(record, idField)
	if err != nil {
		return fmt.Errorf("Unable to determine if record %+v is new - %s", record, err)
	}

	if isNew {
		if err := driver.Create(tx, meta.tablename, fields, &idField); err != nil {
			return fmt.Errorf("Unable to create record %+v - %s", record, err)
		}

		if err := assignField(record, idField); err != nil {
			return fmt.Errorf("Unable to assign primary field for record %+v - %s", record, err)
		}
	} else {
		if err := populateFieldValue(record, &idField); err != nil {
			return fmt.Errorf("Unable to fetch primary field from record %+v - %s", record, err)
		}

		if err := driver.Update(tx, meta.tablename, fields, idField); err != nil {
			return fmt.Errorf("Unable to update record %+v - %s", record, err)
		}
	}

	return nil
}

func remove(tx interface{}, record interface{}) error {
	meta, err := getMetadata(record)
	if err != nil {
		return err
	}

	idField := meta.primary
	if err := populateFieldValue(record, &idField); err != nil {
		return fmt.Errorf("Unable to populate primary field of record %+v - %s", record, err)
	}

	if err := driver.Remove(tx, meta.tablename, idField); err != nil {
		return fmt.Errorf("Unable to remove record %+v - %s", record, err)
	}

	return nil
}
