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
	setup(t)

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

func TestSaveUpdates(t *testing.T) {
	setup(t)

	p := &Person{Name: "John", Age: 31}
	if err := rebecca.Save(p); err != nil {
		t.Fatal(err)
	}

	expected := &Person{}
	if err := rebecca.Get(p.ID, expected); err != nil {
		t.Fatal(err)
	}

	expected.Name = "John Smith"
	if err := rebecca.Save(expected); err != nil {
		t.Fatal(err)
	}

	actual := &Person{}
	if err := rebecca.Get(p.ID, actual); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected %+v to equal %+v", actual, expected)
	}
}

func (d *Driver) exec(t *testing.T, query string) {
	if _, err := d.db.Exec(query); err != nil {
		t.Fatal(err)
	}
}

func setup(t *testing.T) {
	d := NewDriver(pgURL)
	d.exec(t, "DELETE FROM people")
	d.exec(t, "DELETE FROM posts")
	rebecca.SetupDriver(d)
}
