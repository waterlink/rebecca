package rebecca

import (
	"os"
	"reflect"
	"testing"

	"github.com/waterlink/rebecca/fakedriver"
)

func TestMain(m *testing.M) {
	SetupDriver(fakedriver.NewDriver())
	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestSaveCreates(t *testing.T) {
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
