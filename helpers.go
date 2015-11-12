package rebecca

// This file contains shared functions for rebecca package.

import (
	"fmt"
	"reflect"

	"github.com/waterlink/rebecca/driver"
	"github.com/waterlink/rebecca/field"
)

func get(tx interface{}, ID interface{}, record interface{}) error {
	d, lock := driver.Get()
	defer lock.Unlock()

	meta, err := getMetadata(record)
	if err != nil {
		return err
	}

	idField := meta.primary
	idField.Value = ID

	fields, err := d.Get(tx, meta.tablename, meta.fields, idField)
	if err != nil {
		return fmt.Errorf("Unable to find record - %s", err)
	}

	if err := setFields(record, fields); err != nil {
		return fmt.Errorf("Unable to construct found record - %s", err)
	}

	return nil
}

func save(tx interface{}, record interface{}) error {
	d, lock := driver.Get()
	defer lock.Unlock()

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
		if err := d.Create(tx, meta.tablename, fields, &idField); err != nil {
			return fmt.Errorf("Unable to create record %+v - %s", record, err)
		}

		if err := assignField(record, idField); err != nil {
			return fmt.Errorf("Unable to assign primary field for record %+v - %s", record, err)
		}
	} else {
		if err := populateFieldValue(record, &idField); err != nil {
			return fmt.Errorf("Unable to fetch primary field from record %+v - %s", record, err)
		}

		if err := d.Update(tx, meta.tablename, fields, idField); err != nil {
			return fmt.Errorf("Unable to update record %+v - %s", record, err)
		}
	}

	return nil
}

func remove(tx interface{}, record interface{}) error {
	d, lock := driver.Get()
	defer lock.Unlock()

	meta, err := getMetadata(record)
	if err != nil {
		return err
	}

	idField := meta.primary
	if err := populateFieldValue(record, &idField); err != nil {
		return fmt.Errorf("Unable to populate primary field of record %+v - %s", record, err)
	}

	if err := d.Remove(tx, meta.tablename, idField); err != nil {
		return fmt.Errorf("Unable to remove record %+v - %s", record, err)
	}

	return nil
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
