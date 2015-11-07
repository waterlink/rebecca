package rebecca

import (
	"reflect"
	"testing"

	"github.com/waterlink/rebecca/fakedriver"
)

func TestSaveCreates(t *testing.T) {
	SetupDriver(fakedriver.NewDriver())

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
	if err := Get(expected.ID, actual); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %+v to equal %+v", actual, expected)
	}
}

func TestSaveUpdates(t *testing.T) {
	SetupDriver(fakedriver.NewDriver())

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
	if err := Get(p.ID, expected); err != nil {
		t.Fatal(err)
	}

	expected.Name = "John Smith Jr"
	if err := Save(expected); err != nil {
		t.Fatal(err)
	}

	actual := &Person{}
	if err := Get(p.ID, actual); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %+v to equal %+v", actual, expected)
	}
}

func TestAll(t *testing.T) {
	SetupDriver(fakedriver.NewDriver())

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
