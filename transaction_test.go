package rebecca_test

import (
	"errors"
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/waterlink/rebecca"
	"github.com/waterlink/rebecca/driver/fake"
	"github.com/waterlink/rebecca/field"
)

func ExampleTransact() {
	type Person struct {
		// ...
	}

	rebecca.Transact(func(tx *rebecca.Transaction) error {
		// Now you can use `tx` the same way as `rebecca` package, i.e.:
		people := []Person{}
		if err := tx.Where(&people, "name = $1 AND age > $2", "James", 25); err != nil {
			// returning non-nil result here will make transaction roll back
			return err
			// panicking will achieve the same result
			// panic(err)
		}

		// At this point people contains all Person records with name="James" and
		// with age > 25.
		fmt.Print(people)

		// This way you can use all main exported functions of rebecca package as
		// methods on `tx`:
		// - tx.All(records)
		// - tx.First(record, where, args...)
		// - tx.Get(record, ID)
		// - tx.Remove(record)
		// - tx.Save(record)
		// - tx.Where(records, where, args...)

		// For example:
		record := &Person{}
		return tx.Save(record)
	})
}

func ExampleBegin() {
	type Person struct {
		// ...
	}

	// Lets begin transaction:
	tx, err := rebecca.Begin()
	if err != nil {
		// Handle error when unable to begin transaction here.
		panic(err)
	}
	// Lets make sure, that transaction gets rolled back if we return prematurely:
	defer tx.Rollback()

	// Now you can use `tx` the same way as `rebecca` package, i.e.:
	people := []Person{}
	if err := tx.Where(&people, "name = $1 AND age > $2", "James", 25); err != nil {
		panic(err)
	}

	// At this point people contains all Person records with name="James" and
	// with age > 25.
	fmt.Print(people)

	// This way you can use all main exported functions of rebecca package as
	// methods on `tx`:
	// - tx.All(records)
	// - tx.First(record, where, args...)
	// - tx.Get(record, ID)
	// - tx.Remove(record)
	// - tx.Save(record)
	// - tx.Where(records, where, args...)
}

func ExampleTransaction_Commit() {
	tx, err := rebecca.Begin()
	if err != nil {
		panic(err)
	}
	defer tx.Rollback()

	// .. doing some hard work with tx ..

	// And finally, lets commit the transaction:
	if err := tx.Commit(); err != nil {
		// Handle error, when transaction can not be committed, here.
		panic(err)
	}
	// At this point transaction `tx` is committed and should not be used
	// further.
}

func ExampleTransaction_Rollback() {
	tx, err := rebecca.Begin()
	if err != nil {
		panic(err)
	}
	defer tx.Rollback()

	// Sometimes `defer tx.Rollback()` is not acceptable and you might need
	// better control, in that case, you can just call `tx.Rollback()` when
	// necessary:
	if someBadCondition() {
		tx.Rollback()
	}
}

func ExampleTransaction_Context() {
	type Person struct {
		// ...
	}

	tx, err := rebecca.Begin()
	if err != nil {
		panic(err)
	}
	defer tx.Rollback()

	// When you need to use rebecca.Context features (like Order or Limit/Skip,
	// or even Group) together with transaction, you can instantiate
	// rebecca.Context using transaction method Context:
	ctx := tx.Context(&rebecca.Context{Order: "age ASC", Limit: 30, Skip: 90})

	// And then use `ctx` as usual:
	people := []Person{}
	if err := ctx.All(&people); err != nil {
		panic(err)
	}
	fmt.Print(people)
}

func someBadCondition() bool {
	return true
}

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

func TestTransactions(t *testing.T) {
	d := fake.NewDriver()
	rebecca.SetupDriver(d)

	d.RegisterWhere("title = $1", func(record []field.Field, args ...interface{}) (bool, error) {
		for _, f := range record {
			if f.DriverName == "title" {
				return f.Value.(string) == args[0].(string), nil
			}
		}

		return false, fmt.Errorf("record %+v does not have title field", record)
	})

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

	pa2 := &Post{Title: "Hello Blog", Content: "More hellos here!", CreatedAt: time.Now()}
	if err := txa.Save(pa2); err != nil {
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

	expecteds := []Post{*pa, *pa2}
	actuals := []Post{}
	if err := txa.All(&actuals); err != nil {
		t.Fatal(err)
	}

	if !equalPosts(actuals, expecteds) {
		t.Errorf("Expected %+v to equal %+v", actuals, expecteds)
	}

	expecteds = []Post{*pa}
	actuals = []Post{}
	if err := txa.Where(&actuals, "title = $1", "Hello world"); err != nil {
		t.Fatal(err)
	}

	if !equalPosts(actuals, expecteds) {
		t.Errorf("Expected %+v to equal %+v", actuals, expecteds)
	}

	expecteds = []Post{*pa2}
	actuals = []Post{}
	if err := txa.Where(&actuals, "title = $1", "Hello Blog"); err != nil {
		t.Fatal(err)
	}

	if !equalPosts(actuals, expecteds) {
		t.Errorf("Expected %+v to equal %+v", actuals, expecteds)
	}

	expecteds = []Post{}
	actuals = []Post{}
	if err := txa.Where(&actuals, "title = $1", "Super Post"); err != nil {
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

	expecteds = []Post{}
	actuals = []Post{}
	if err := txc.Where(&actuals, "title = $1", "Super Post"); err != nil {
		t.Fatal(err)
	}

	if !equalPosts(actuals, expecteds) {
		t.Errorf("Expected %+v to equal %+v", actuals, expecteds)
	}

	expected := pa2
	actual = &Post{}
	if err := txc.First(actual, "title = $1", "Hello Blog"); err != nil {
		t.Fatal(err)
	}

	if !actual.Equal(expected) {
		t.Errorf("Expected %+v to equal %+v", actual, expected)
	}

	txd, err := rebecca.Begin()
	if err != nil {
		t.Fatal(err)
	}

	if err := txd.Remove(pa); err != nil {
		t.Fatal(err)
	}

	expecteds = []Post{*pa2}
	actuals = []Post{}
	if err := txd.All(&actuals); err != nil {
		t.Fatal(err)
	}

	if !equalPosts(actuals, expecteds) {
		t.Errorf("Expected %+v to equal %+v", actuals, expecteds)
	}

	expecteds = []Post{*pa, *pa2}
	actuals = []Post{}
	if err := txc.All(&actuals); err != nil {
		t.Fatal(err)
	}

	if !equalPosts(actuals, expecteds) {
		t.Errorf("Expected %+v to equal %+v", actuals, expecteds)
	}

	expecteds = []Post{*pa, *pa2}
	actuals = []Post{}
	if err := rebecca.All(&actuals); err != nil {
		t.Fatal(err)
	}

	if !equalPosts(actuals, expecteds) {
		t.Errorf("Expected %+v to equal %+v", actuals, expecteds)
	}

	if err := txd.Commit(); err != nil {
		t.Fatal(err)
	}

	expecteds = []Post{*pa2}
	actuals = []Post{}
	if err := rebecca.All(&actuals); err != nil {
		t.Fatal(err)
	}

	if !equalPosts(actuals, expecteds) {
		t.Errorf("Expected %+v to equal %+v", actuals, expecteds)
	}
}

func TestDoubleCommit(t *testing.T) {
	rebecca.SetupDriver(fake.NewDriver())

	tx, err := rebecca.Begin()
	if err != nil {
		t.Fatal(err)
	}

	if err := tx.Commit(); err != nil {
		t.Fatal(err)
	}

	err = tx.Commit()
	if err == nil {
		t.Fatal("Expected transaction to no being able to be committed twice")
	}

	expected := `Unable to commit transaction - Current transaction is already finished`
	actual := err.Error()
	if actual != expected {
		t.Errorf("Expected %s to equal %s", actual, expected)
	}
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

func TestTransact(t *testing.T) {
	now := time.Now()
	var p *Post

	examples := map[string]struct {
		handler func(tx *rebecca.Transaction) error
		verify  func(err error) error
	}{

		"when handler succeeds": {
			handler: func(tx *rebecca.Transaction) error {
				p = &Post{Title: "Hello", Content: "World", CreatedAt: now}
				return tx.Save(p)
			},

			verify: func(err error) error {
				if err != nil {
					return fmt.Errorf("Expected handler to not fail, but got: %s", err)
				}

				actual := &Post{}
				if err := rebecca.Get(actual, p.ID); err != nil {
					return fmt.Errorf("Expected rebecca.Get(actual, %d) to not fail, but got: %s", p.ID, err)
				}

				if !actual.Equal(p) {
					return fmt.Errorf("Expected %#v to equal %#v", actual, p)
				}

				return nil
			},
		},

		"when handler fails": {
			handler: func(tx *rebecca.Transaction) error {
				p = &Post{Title: "Hello", Content: "World", CreatedAt: now}
				tx.Save(p)
				return errors.New("I have failed")
			},

			verify: func(err error) error {
				if err == nil {
					return errors.New("Expected handler to fail, but got nil")
				}

				if err.Error() != "I have failed" {
					return fmt.Errorf("Expected error to be 'I have failed', but got: '%s'", err)
				}

				actual := &Post{}
				if err := rebecca.Get(actual, p.ID); err == nil {
					return fmt.Errorf("Expected rebecca.Get(actual, %d) to fail, but got nil", p.ID)
				}

				return nil
			},
		},

		"when handler panics": {
			handler: func(tx *rebecca.Transaction) error {
				p = &Post{Title: "Hello", Content: "World", CreatedAt: now}
				tx.Save(p)
				panic("I have a panic!")
			},

			verify: func(err error) error {
				if err == nil {
					return errors.New("Expected handler to fail, but got nil")
				}

				if err.Error() != "I have a panic! (recovered)" {
					return fmt.Errorf("Expected error to be 'I have a panic! (recovered)', but got: '%s'", err)
				}

				actual := &Post{}
				if err := rebecca.Get(actual, p.ID); err == nil {
					return fmt.Errorf("Expected rebecca.Get(actual, %d) to fail, but got nil", p.ID)
				}

				return nil
			},
		},
	}

	for info, e := range examples {
		t.Log(info)
		p = nil
		if err := e.verify(rebecca.Transact(e.handler)); err != nil {
			t.Error(err)
		}
	}
}
