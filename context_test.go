package rebecca_test

import (
	"fmt"

	"github.com/waterlink/rebecca"
)

func ExampleContext_All() {
	type Person struct {
		// ...
	}

	ctx := rebecca.Context{Limit: 20, Skip: 40}
	people := []Person{}
	if err := ctx.All(&people); err != nil {
		panic(err)
	}
	// At this point people contains 20 Person records starting from 41th from
	// the database.
	fmt.Print(people)
}

func ExampleContext_First() {
	type Person struct {
		// ...
	}

	ctx := rebecca.Context{Order: "age DESC"}
	oldestTeenager := &Person{}
	if err := ctx.First(oldestTeenager, "age < $1", 21); err != nil {
		panic(err)
	}
	// At this point oldestTeenager will contain a Person record that is of
	// maximum age and that is of age < 21.
	fmt.Print(oldestTeenager)
}

func ExampleContext_Where() {
	type Person struct {
		// ...
	}

	ctx := rebecca.Context{Order: "age DESC"}
	teenagers := []Person{}
	if err := ctx.Where(&teenagers, "age < $1", 21); err != nil {
		panic(err)
	}
	// At this point teenagers will contain a list of Person records sorted by
	// age in descending order and where age < 21.
	fmt.Print(teenagers)
}
