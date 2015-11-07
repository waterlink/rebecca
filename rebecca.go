package rebecca

import (
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
	Get(tablename string, fields []field.Field, ID field.Field) ([]field.Field, error)
	Create(tablename string, fields []field.Field, ID *field.Field) error
	Update(tablename string, fields []field.Field, ID field.Field) error
	All(tablename string, fields []field.Field, ctx *context.Context) ([][]field.Field, error)
	Where(tablename string, fields []field.Field, ctx *context.Context, where string) ([][]field.Field, error)
	First(tablename string, fields []field.Field, ctx *context.Context, where string) ([]field.Field, error)
	Remove(tablename string, ID field.Field) error
}

// Context is for storing query context
type Context struct {
	Order string
	Group string
	Limit int
	Skip  int
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

// All is for fetching all records
func (c *Context) All(records []interface{}) error {
	return nil
}

// Where is for fetching specific records
func (c *Context) Where(query string, records []interface{}) error {
	return nil
}

// First is for fetching only one specific record
func (c *Context) First(query string, record interface{}) error {
	return nil
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
	meta, err := getMetadata(record)
	if err != nil {
		return err
	}

	idField := meta.primary
	idField.Value = ID

	fields, err := driver.Get(meta.tablename, meta.fields, idField)
	if err != nil {
		return fmt.Errorf("Unable to find record - %s", err)
	}

	if err := setFields(record, fields); err != nil {
		return fmt.Errorf("Unable to construct found record - %s", err)
	}

	return nil
}

// Save is for saving one record (either creating or updating)
func Save(record interface{}) error {
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
		if err := driver.Create(meta.tablename, fields, &idField); err != nil {
			return fmt.Errorf("Unable to create record %+v - %s", record, err)
		}

		if err := assignField(record, idField); err != nil {
			return fmt.Errorf("Unable to assign primary field for record %+v - %s", record, err)
		}
	} else {
		if err := populateFieldValue(record, &idField); err != nil {
			return fmt.Errorf("Unable to fetch primary field from record %+v - %s", record, err)
		}

		if err := driver.Update(meta.tablename, fields, idField); err != nil {
			return fmt.Errorf("Unable to update record %+v - %s", record, err)
		}
	}

	return nil
}

// All is for fetching all records
func All(records []interface{}) error {
	ctx := &Context{}
	return ctx.All(records)
}

// Where is for fetching specific records
func Where(where string, records []interface{}) error {
	ctx := &Context{}
	return ctx.Where(where, records)
}

// First is for fetching only one specific record
func First(where string, record interface{}) error {
	ctx := &Context{}
	return ctx.First(where, record)
}

// Remove is for removing the record
func Remove(record interface{}) error {
	return nil
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

func fetchMetadata(record interface{}) (metadata, error) {
	missingMetadata := metadata{}

	v := reflect.ValueOf(record)
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return missingMetadata, fmt.Errorf("Rebecca's model is required to be struct, but got: %+v", record)
	}

	metaValue := v.FieldByName("ModelMetadata")
	metaType := metaValue.Type()
	if metaType.PkgPath() != "github.com/waterlink/rebecca" || metaType.Name() != "ModelMetadata" {
		return missingMetadata, fmt.Errorf("Rebecca's model is required to embed rebecca.ModelMetadata")
	}

	meta := metadata{}

	ty := v.Type()
	metaField, ok := ty.FieldByName("ModelMetadata")
	if !ok {
		return missingMetadata, fmt.Errorf("Unable to find ModelMetadata field on %s.%s type", ty.PkgPath(), ty.Name())
	}

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
	for v.Kind() == reflect.Ptr {
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
	for v.Kind() == reflect.Ptr {
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
	for v.Kind() == reflect.Ptr {
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
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	ty := v.Type()

	if _, ok := ty.FieldByName(ID.Name); !ok {
		return false, fmt.Errorf("Field %s not found on record %+v", ID.Name, record)
	}

	f := v.FieldByName(ID.Name)
	return f.Interface() == reflect.Zero(f.Type()).Interface(), nil
}
