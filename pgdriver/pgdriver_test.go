package pgdriver

import (
	"reflect"
	"testing"
	"time"

	"github.com/waterlink/rebecca"
)

const (
	pgURL = "postgres://rebecca_pg:rebecca_pg@127.0.0.1:5432/rebecca_pg_test?sslmode=disable"
)

type Person struct {
	rebecca.ModelMetadata `tablename:"people"`

	ID   int    `rebecca:"id" rebecca_primary:"true"`
	Name string `rebecca:"name"`
	Age  int    `rebecca:"age"`
}

type Post struct {
	rebecca.ModelMetadata `tablename:"posts"`

	ID        int       `rebecca;"id" rebecca_primary:"true"`
	Title     string    `rebecca:"title"`
	Content   string    `rebecca:"content"`
	CreatedAt time.Time `rebecca:"created_at"`
}

func TestSaveCreates(t *testing.T) {
	rebecca.SetupDriver(NewDriver(pgURL))

	expected := &Person{Name: "John", Age: 31}
	if err := rebecca.Save(expected); err != nil {
		t.Fatal(err)
	}

	actual := &Person{}
	if err := rebecca.Get(expected.ID, actual); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected %+v to equal %+v", actual, expected)
	}
}
