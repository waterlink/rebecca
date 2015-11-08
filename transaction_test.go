package rebecca_test

import (
	"fmt"

	"github.com/waterlink/rebecca"
)

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
	// - tx.Get(ID, record)
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
