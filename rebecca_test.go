package rebecca

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/waterlink/rebecca/driver/fake"
	"github.com/waterlink/rebecca/field"
)

func TestSaveCreates(t *testing.T) {
	SetupDriver(fake.NewDriver())

	type Person struct {
		ModelMetadata `tablename:"people"`

		ID   int    `rebecca:"id" rebecca_primary:"true"`
		Name string `rebecca:"name"`
		Age  int    `rebecca:"age"`
	}

	expected := &Person{Name: "John Smith", Age: 31}
	if err := Save(expected); err != nil {
		t.Fatal(err)
	}

	actual := &Person{}
	if err := Get(actual, expected.ID); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %+v to equal %+v", actual, expected)
	}
}

func TestSaveUpdates(t *testing.T) {
	SetupDriver(fake.NewDriver())

	type Person struct {
		ModelMetadata `tablename:"people"`

		ID   int    `rebecca:"id" rebecca_primary:"true"`
		Name string `rebecca:"name"`
		Age  int    `rebecca:"age"`
	}

	p := &Person{Name: "John Smith", Age: 31}
	if err := Save(p); err != nil {
		t.Fatal(err)
	}

	expected := &Person{}
	if err := Get(expected, p.ID); err != nil {
		t.Fatal(err)
	}

	expected.Name = "John Smith Jr"
	if err := Save(expected); err != nil {
		t.Fatal(err)
	}

	actual := &Person{}
	if err := Get(actual, p.ID); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %+v to equal %+v", actual, expected)
	}
}

func TestAll(t *testing.T) {
	SetupDriver(fake.NewDriver())

	type Person struct {
		ModelMetadata `tablename:"people"`

		ID   int    `rebecca:"id" rebecca_primary:"true"`
		Name string `rebecca:"name"`
		Age  int    `rebecca:"age"`
	}

	p1 := &Person{Name: "John", Age: 37}
	p2 := &Person{Name: "Sarah", Age: 26}
	p3 := &Person{Name: "James", Age: 33}
	people := []*Person{p1, p2, p3}

	for _, p := range people {
		if err := Save(p); err != nil {
			t.Fatal(err)
		}
	}

	expected := []Person{*p1, *p2, *p3}
	actual := []Person{}
	if err := All(&actual); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected %+v to equal to %+v", actual, expected)
	}
}

func TestWhere(t *testing.T) {
	d := fake.NewDriver()
	SetupDriver(d)

	d.RegisterWhere("age < $1", func(record []field.Field, args ...interface{}) (bool, error) {
		for _, f := range record {
			if f.DriverName == "age" {
				return f.Value.(int) < args[0].(int), nil
			}
		}

		return false, fmt.Errorf("record %+v does not have age field", record)
	})

	d.RegisterWhere("age >= $1", func(record []field.Field, args ...interface{}) (bool, error) {
		for _, f := range record {
			if f.DriverName == "age" {
				return f.Value.(int) >= args[0].(int), nil
			}
		}

		return false, fmt.Errorf("record %+v does not have age field", record)
	})

	type Person struct {
		ModelMetadata `tablename:"people"`

		ID   int    `rebecca:"id" rebecca_primary:"true"`
		Name string `rebecca:"name"`
		Age  int    `rebecca:"age"`
	}

	p1 := &Person{Name: "John", Age: 9}
	p2 := &Person{Name: "Sarah", Age: 27}
	p3 := &Person{Name: "James", Age: 11}
	p4 := &Person{Name: "Monika", Age: 12}
	people := []*Person{p1, p2, p3, p4}

	for _, p := range people {
		if err := Save(p); err != nil {
			t.Fatal(err)
		}
	}

	expected := []Person{*p1, *p3}
	actual := []Person{}
	if err := Where(&actual, "age < $1", 12); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected %+v to equal to %+v", actual, expected)
	}

	expected = []Person{*p2, *p4}
	actual = []Person{}
	if err := Where(&actual, "age >= $1", 12); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected %+v to equal to %+v", actual, expected)
	}
}

func TestFirst(t *testing.T) {
	d := fake.NewDriver()
	SetupDriver(d)

	d.RegisterWhere("age < $1", func(record []field.Field, args ...interface{}) (bool, error) {
		for _, f := range record {
			if f.DriverName == "age" {
				return f.Value.(int) < args[0].(int), nil
			}
		}

		return false, fmt.Errorf("record %+v does not have age field", record)
	})

	d.RegisterWhere("age >= $1", func(record []field.Field, args ...interface{}) (bool, error) {
		for _, f := range record {
			if f.DriverName == "age" {
				return f.Value.(int) >= args[0].(int), nil
			}
		}

		return false, fmt.Errorf("record %+v does not have age field", record)
	})

	type Person struct {
		ModelMetadata `tablename:"people"`

		ID   int    `rebecca:"id" rebecca_primary:"true"`
		Name string `rebecca:"name"`
		Age  int    `rebecca:"age"`
	}

	p1 := &Person{Name: "John", Age: 9}
	p2 := &Person{Name: "Sarah", Age: 27}
	p3 := &Person{Name: "James", Age: 11}
	p4 := &Person{Name: "Monika", Age: 12}
	people := []*Person{p1, p2, p3, p4}

	for _, p := range people {
		if err := Save(p); err != nil {
			t.Fatal(err)
		}
	}

	expected := p1
	actual := &Person{}
	if err := First(&actual, "age < $1", 12); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected %+v to equal to %+v", actual, expected)
	}

	expected = p2
	actual = &Person{}
	if err := First(&actual, "age >= $1", 12); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected %+v to equal to %+v", actual, expected)
	}
}

func TestRemove(t *testing.T) {
	d := fake.NewDriver()
	SetupDriver(d)

	type Person struct {
		ModelMetadata `tablename:"people"`

		ID   int    `rebecca:"id" rebecca_primary:"true"`
		Name string `rebecca:"name"`
		Age  int    `rebecca:"age"`
	}

	p1 := &Person{Name: "John", Age: 9}
	p2 := &Person{Name: "Sarah", Age: 27}
	p3 := &Person{Name: "James", Age: 11}
	p4 := &Person{Name: "Monika", Age: 12}
	people := []*Person{p1, p2, p3, p4}

	for _, p := range people {
		if err := Save(p); err != nil {
			t.Fatal(err)
		}
	}

	if err := Remove(p2); err != nil {
		t.Fatal(err)
	}

	expected := []Person{*p1, *p3, *p4}
	actual := []Person{}
	if err := All(&actual); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected %+v to equal to %+v", actual, expected)
	}
}

func TestInvalidMetadata(t *testing.T) {
	type Person struct {
		ModelMetadata `tablename:"people"`

		ID   int    `rebecca:"id" rebecca_primary:"true"`
		Name string `rebecca:"name"`
	}

	type MissingMetadata struct {
		ID   int    `rebecca:"id" rebecca_primary:"true"`
		Name string `rebecca:"name"`
	}

	type IncorrectMetadata struct {
		ModelMetadata string `tablename:"people"`
		ID            int    `rebecca:"id" rebecca_primary:"true"`
		Name          string `rebecca:"name"`
	}

	type MissingTablename struct {
		ModelMetadata
		ID   int    `rebecca:"id" rebecca_primary:"true"`
		Name string `rebecca:"name"`
	}

	type NoPrimary struct {
		ModelMetadata `tablename:"noprimaries"`

		ID   int    `rebecca:"id"`
		Name string `rebecca:"name"`
	}

	type NonStruct []int

	missingMetadataError := `Unable to fetch record's metadata - type=github.com/waterlink/rebecca.MissingMetadata - Rebecca's model is required to embed rebecca.ModelMetadata`
	missingPrimaryField := "Record has no primary field - type=github.com/waterlink/rebecca.NoPrimary - Use `rebecca_primary:\"true\"` annotation"
	nonStructError := `Unable to fetch record's metadata - type=.int - Rebecca's model is required to be struct, but got: &[]`
	incorrectMetadataError := `Unable to fetch record's metadata - type=github.com/waterlink/rebecca.IncorrectMetadata - Rebecca's model is required to embed rebecca.ModelMetadata`
	missingTablenameError := `Unable to fetch record's metadata - type=github.com/waterlink/rebecca.MissingTablename - tablename tag metadata is missing on ModelMetadata`
	notAccessibleError := `Unable to assign primary field for record {ModelMetadata:{} ID:0 Name:} - Unable to set field ID on record {ModelMetadata:{} ID:0 Name:}. It is required to be exported and addressable`
	notAccessibleOnGet := `Unable to construct found record - Unable to set field ID on record {ModelMetadata:{} ID:0 Name:}. It is required to be exported and addressable`

	examples := map[string]struct {
		thing  interface{}
		action func(interface{}) error
		err    error
	}{
		"Save fails when metadata is missing": {
			thing:  &MissingMetadata{Name: "stuff"},
			action: func(x interface{}) error { return Save(x) },
			err:    errors.New(missingMetadataError),
		},

		"Save fails when metadata's tablename is missing": {
			thing:  &MissingTablename{Name: "stuff"},
			action: func(x interface{}) error { return Save(x) },
			err:    errors.New(missingTablenameError),
		},

		"Save fails when metadata is of incorrect type": {
			thing:  &IncorrectMetadata{Name: "stuff"},
			action: func(x interface{}) error { return Save(x) },
			err:    errors.New(incorrectMetadataError),
		},

		"Save fails when there is no primary key defined": {
			thing:  &NoPrimary{Name: "James"},
			action: func(x interface{}) error { return Save(x) },
			err:    errors.New(missingPrimaryField),
		},

		"Save fails when model is not struct": {
			thing:  &NonStruct{},
			action: func(x interface{}) error { return Save(x) },
			err:    errors.New(nonStructError),
		},

		"Save fails when thing is not accessible": {
			thing:  Person{},
			action: func(x interface{}) error { return Save(x) },
			err:    errors.New(notAccessibleError),
		},

		"Get fails when metadata is missing": {
			thing:  &MissingMetadata{},
			action: func(x interface{}) error { return Get(x, 123) },
			err:    errors.New(missingMetadataError),
		},

		"Get fails when metadata is of incorrect type": {
			thing:  &IncorrectMetadata{},
			action: func(x interface{}) error { return Get(x, 123) },
			err:    errors.New(incorrectMetadataError),
		},

		"Get fails when model is not struct": {
			thing:  &NonStruct{},
			action: func(x interface{}) error { return Get(x, "hello") },
			err:    errors.New(nonStructError),
		},

		"Get fails when thing is not accessible": {
			thing: Person{},
			action: func(x interface{}) error {
				p := x.(Person)
				Save(&p)
				return Get(x, p.ID)
			},
			err: errors.New(notAccessibleOnGet),
		},

		"Remove fails when metadata is missing": {
			thing:  &MissingMetadata{Name: "John", ID: 123},
			action: func(x interface{}) error { return Remove(x) },
			err:    errors.New(missingMetadataError),
		},

		"Remove fails when metadata is of incorrect type": {
			thing:  &IncorrectMetadata{Name: "John", ID: 123},
			action: func(x interface{}) error { return Remove(x) },
			err:    errors.New(incorrectMetadataError),
		},

		"Remove fails when there is no primary key defined": {
			thing:  &NoPrimary{Name: "James", ID: 123},
			action: func(x interface{}) error { return Remove(x) },
			err:    errors.New(missingPrimaryField),
		},

		"Remove fails when model is not struct": {
			thing:  &NonStruct{},
			action: func(x interface{}) error { return Remove(x) },
			err:    errors.New(nonStructError),
		},
	}

	for info, e := range examples {
		t.Log(info)
		err := e.action(e.thing)
		actual := errRepr(err)
		expected := errRepr(e.err)
		if actual != expected {
			t.Errorf("Expected %s to equal %s", actual, expected)
		}
	}
}

func errRepr(err error) string {
	if err == nil {
		return "<nil :: error>"
	}
	return err.Error()
}
