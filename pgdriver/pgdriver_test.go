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

func TestFirst(t *testing.T) {
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

	expected := p2
	actual := &Person{}
	if err := rebecca.First(actual, "age > $1", 10); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected %+v to equal %+v", actual, expected)
	}

	expected = p1
	actual = &Person{}
	if err := rebecca.First(actual, "age < $1", 15); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected %+v to equal %+v", actual, expected)
	}
}

func TestLimit(t *testing.T) {
	setup(t)

	p1 := &Person{Name: "John", Age: 9}
	p2 := &Person{Name: "Sarah", Age: 27}
	p3 := &Person{Name: "Bruce", Age: 33}
	p4 := &Person{Name: "James", Age: 11}
	p5 := &Person{Name: "Monika", Age: 12}
	p6 := &Person{Name: "Peter", Age: 21}
	people := []*Person{p1, p2, p3, p4, p5, p6}

	for _, p := range people {
		if err := rebecca.Save(p); err != nil {
			t.Fatal(err)
		}
	}

	ctx := &rebecca.Context{Limit: 4}
	expected := []Person{*p1, *p2, *p3, *p4}
	actual := []Person{}
	if err := ctx.All(&actual); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected %+v to equal %+v", actual, expected)
	}

	expected = []Person{*p2, *p3, *p4, *p5}
	actual = []Person{}
	if err := ctx.Where(&actual, "age > $1", 10); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected %+v to equal %+v", actual, expected)
	}
}

func TestSkip(t *testing.T) {
	setup(t)

	p1 := &Person{Name: "John", Age: 9}
	p2 := &Person{Name: "Sarah", Age: 27}
	p3 := &Person{Name: "Bruce", Age: 33}
	p4 := &Person{Name: "James", Age: 11}
	p5 := &Person{Name: "Monika", Age: 12}
	p6 := &Person{Name: "Peter", Age: 21}
	people := []*Person{p1, p2, p3, p4, p5, p6}

	for _, p := range people {
		if err := rebecca.Save(p); err != nil {
			t.Fatal(err)
		}
	}

	ctx := &rebecca.Context{Skip: 2}
	expected := []Person{*p3, *p4, *p5, *p6}
	actual := []Person{}
	if err := ctx.All(&actual); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected %+v to equal %+v", actual, expected)
	}

	expected = []Person{*p5, *p6}
	actual = []Person{}
	if err := ctx.Where(&actual, "age > $1", 11); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected %+v to equal %+v", actual, expected)
	}

	expectedOne := p4
	actualOne := &Person{}
	if err := ctx.First(&actualOne, "age > $1", 10); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(actualOne, expectedOne) {
		t.Errorf("Expected %+v to equal %+v", actualOne, expectedOne)
	}
}

func TestOrder(t *testing.T) {
	setup(t)

	p1 := &Person{Name: "John", Age: 9}
	p2 := &Person{Name: "Sarah", Age: 27}
	p3 := &Person{Name: "Bruce", Age: 33}
	p4 := &Person{Name: "James", Age: 11}
	p5 := &Person{Name: "Monika", Age: 12}
	p6 := &Person{Name: "Peter", Age: 21}
	people := []*Person{p1, p2, p3, p4, p5, p6}

	for _, p := range people {
		if err := rebecca.Save(p); err != nil {
			t.Fatal(err)
		}
	}

	ctx := &rebecca.Context{Order: "age DESC"}
	expected := []Person{*p3, *p2, *p6, *p5, *p4, *p1}
	actual := []Person{}
	if err := ctx.All(&actual); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected %+v to equal %+v", actual, expected)
	}

	expected = []Person{*p6, *p5, *p4, *p1}
	actual = []Person{}
	if err := ctx.Where(&actual, "age <= $1", 21); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected %+v to equal %+v", actual, expected)
	}

	expectedOne := p5
	actualOne := &Person{}
	if err := ctx.First(actualOne, "age < $1", 21); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(actualOne, expectedOne) {
		t.Errorf("Expected %+v to equal %+v", actualOne, expectedOne)
	}
}

func TestGroup(t *testing.T) {
	setup(t)

	p1 := &Person{Name: "John", Age: 9}
	p2 := &Person{Name: "Sarah", Age: 27}
	p3 := &Person{Name: "Bruce", Age: 27}
	p4 := &Person{Name: "James", Age: 11}
	p5 := &Person{Name: "Monika", Age: 11}
	p6 := &Person{Name: "Peter", Age: 21}
	p7 := &Person{Name: "Brad", Age: 11}
	people := []*Person{p1, p2, p3, p4, p5, p6, p7}

	for _, p := range people {
		if err := rebecca.Save(p); err != nil {
			t.Fatal(err)
		}
	}

	type PersonByAge struct {
		rebecca.ModelMetadata `tablename:"people"`

		Age   int `rebecca:"age" rebecca_primary:"true"`
		Count int `rebecca:"count(distinct(id))"`
	}

	ctx := rebecca.Context{Group: "age"}
	expected := []PersonByAge{
		{Age: 9, Count: 1},
		{Age: 11, Count: 3},
		{Age: 21, Count: 1},
		{Age: 27, Count: 2},
	}
	actual := []PersonByAge{}
	if err := ctx.All(&actual); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected %+v to equal %+v", actual, expected)
	}
}

func TestRemove(t *testing.T) {
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

	if err := rebecca.Remove(p3); err != nil {
		t.Fatal(err)
	}

	expected := []Person{*p1, *p2, *p4}
	actual := []Person{}
	if err := rebecca.All(&actual); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected %+v to equal %+v", actual, expected)
	}
}

func TestTransactions(t *testing.T) {
	setup(t)

	txa, err := rebecca.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer txa.Rollback()

	txb, err := rebecca.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer txb.Rollback()

	pa := &Post{Title: "Hello world", Content: "Content of Hello World, many hellos", CreatedAt: time.Now()}
	if err := txa.Save(pa); err != nil {
		t.Fatal(err)
	}

	actual := &Post{}
	if err := txa.Get(pa.ID, actual); err != nil {
		t.Fatal(err)
	}

	actual = &Post{}
	if err := txb.Get(pa.ID, actual); err == nil {
		t.Errorf(
			"Expected transaction B not to find record saved in transaction A, but got: %+v",
			actual,
		)
	}

	pb := &Post{Title: "Super Post", Content: "Super Content", CreatedAt: time.Now()}
	if err := txb.Save(pb); err != nil {
		t.Fatal(err)
	}

	actual = &Post{}
	if err := txa.Get(pb.ID, actual); err == nil {
		t.Errorf(
			"Expected transaction A not to find record saved in transaction B, but got: %+v",
			actual,
		)
	}

	expecteds := []Post{*pa}
	actuals := []Post{}
	if err := txa.All(&actuals); err != nil {
		t.Fatal(err)
	}

	if !equalPosts(actuals, expecteds) {
		t.Errorf("Expected %+v to equal %+v", actuals, expecteds)
	}

	expecteds = []Post{*pb}
	actuals = []Post{}
	if err := txb.All(&actuals); err != nil {
		t.Fatal(err)
	}

	if !equalPosts(actuals, expecteds) {
		t.Errorf("Expected %+v to equal %+v", actuals, expecteds)
	}

	if err := txa.Commit(); err != nil {
		t.Fatal(err)
	}

	txb.Rollback()

	txc, err := rebecca.Begin()
	if err != nil {
		t.Fatal(err)
	}

	actual = &Post{}
	if err := txc.Get(pa.ID, actual); err != nil {
		t.Fatal(err)
	}

	if !actual.Equal(pa) {
		t.Errorf("Expected %+v to equal %+v", actual, pa)
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
