package rebecca_test

import (
	"fmt"

	"github.com/waterlink/rebecca"
)

func ExampleAll() {
	type Person struct {
		// ...
	}

	people := []Person{}
	if err := rebecca.All(&people); err != nil {
		panic(err)
	}
	// At this point people contains all Person records.
	fmt.Print(people)
}

func ExampleFirst() {
	type Person struct {
		// ...
	}

	person := &Person{}
	if err := rebecca.First(person, "name = $1", "John Smith"); err != nil {
		panic(err)
	}
	// At this point person contains first record from the database that has
	// name="John Smith".
	fmt.Print(person)
}

func ExampleGet() {
	type Person struct {
		// ...
	}

	person := &Person{}
	if err := rebecca.Get(person, 25); err != nil {
		panic(err)
	}
	// At this point person contains record with primary key equal to 25.
	fmt.Print(person)
}

func ExampleRemove() {
	type Person struct {
		// ...
	}

	// First lets find person with primary key = 25
	person := &Person{}
	if err := rebecca.Get(person, 25); err != nil {
		panic(err)
	}

	// And then, lets remove it
	if err := rebecca.Remove(person); err != nil {
		panic(err)
	}

	// At this point person with primary key = 25 was removed from database.
}

func ExampleSave() {
	// Lets first define our model:
	type Person struct {
		// the table name is people
		rebecca.ModelMetadata `tablename:"people"`

		// ID is of type int, in database it is mapped to `id` and it is a primary
		// key.
		// Name is of type string, in database it is mapped to `name`.
		// Age is of type int, in database it is mapped to `age`.
		ID   int    `rebecca:"id" rebecca_primary:"true"`
		Name string `rebecca:"name"`
		Age  int    `rebecca:"age"`
	}

	// And now, lets create new record:
	person := &Person{Name: "John", Age: 28}
	// And finally, lets save it to the database
	if err := rebecca.Save(person); err != nil {
		panic(err)
	}

	// At this point, record was saved to database as a new record and its ID
	// field was updated accordingly.
	fmt.Print(person)

	// Now lets modify our record:
	person.Name = "John Smith"
	// And save it
	if err := rebecca.Save(person); err != nil {
		panic(err)
	}

	// At this point, original record was update with new data.
	fmt.Print(person)
}

func ExampleWhere() {
	type Person struct {
		// ...
	}

	teenagers := []Person{}
	if err := rebecca.Where(teenagers, "age < $1", 21); err != nil {
		panic(err)
	}

	// At this point teenagers contains all Person records with age < 21.
	fmt.Print(teenagers)
}

func ExampleExec() {
	ID := 25
	if err := rebecca.Exec("UPDATE counters SET value = value + 1 WHERE id = $1", ID); err != nil {
		panic(err)
	}
}
