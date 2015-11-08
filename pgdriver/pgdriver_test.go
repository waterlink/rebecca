package pgdriver

import (
	"math"
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

	ID        int       `rebecca:"id" rebecca_primary:"true"`
	Title     string    `rebecca:"title"`
	Content   string    `rebecca:"content"`
	CreatedAt time.Time `rebecca:"created_at"`
}

func (p *Post) Equal(other *Post) bool {
	return p.ID == other.ID &&
		p.Title == other.Title &&
		p.Content == other.Content &&
		math.Abs(float64(p.CreatedAt.Sub(other.CreatedAt))) < 500 // microseconds
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

func TestAll(t *testing.T) {
	setup(t)

	p1 := &Post{Title: "Post 1", Content: "Content 1", CreatedAt: time.Now()}
	p2 := &Post{Title: "Post 2", Content: "Content 2", CreatedAt: time.Now()}
	p3 := &Post{Title: "Post 3", Content: "Content 3", CreatedAt: time.Now()}
	p4 := &Post{Title: "Post 4", Content: "Content 4", CreatedAt: time.Now()}
	posts := []*Post{p1, p2, p3, p4}

	for _, p := range posts {
		if err := rebecca.Save(p); err != nil {
			t.Fatal(err)
		}
	}

	expected := []Post{*p1, *p2, *p3, *p4}
	actual := []Post{}
	if err := rebecca.All(&actual); err != nil {
		t.Fatal(err)
	}

	if !equalPosts(actual, expected) {
		t.Errorf("Expected %+v to equal %+v", actual, expected)
	}
}

func TestWhere(t *testing.T) {
	setup(t)

	p1 := &Person{Name: "John", Age: 9}
	p2 := &Person{Name: "Sarah", Age: 27}
	p3 := &Person{Name: "James", Age: 11}
	p4 := &Person{Name: "Monika", Age: 12}
	people := []*Person{p1, p2, p3, p4}

	for _, p := range people {
		if err := rebecca.Save(p); err != nil {
			t.Fatal(err)
		}
	}

	expected := []Person{*p1, *p3}
	actual := []Person{}
	if err := rebecca.Where(&actual, "age < $1", 12); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected %+v to equal %+v", actual, expected)
	}

	expected = []Person{*p2, *p4}
	actual = []Person{}
	if err := rebecca.Where(&actual, "age > $1", 11); err != nil {
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

func equalPosts(l, r []Post) bool {
	if len(l) != len(r) {
		return false
	}

	for i := range l {
		if !l[i].Equal(&r[i]) {
			return false
		}
	}

	return true
}
