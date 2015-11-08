package rebecca_test

import (
	"fmt"

	"github.com/waterlink/rebecca"
)

func ExampleModelMetadata() {
	type Person struct {
		// ModelMetadata allows to provide table name for the model
		rebecca.ModelMetadata `tablename:"people"`

		ID   int    `rebecca:"id" rebecca_primary:"true"`
		Name string `rebecca:"name"`
		Age  int    `rebecca:"age"`
	}

	type PostsOfPerson struct {
		// Additionally you can have any expression as a table name, that your
		// driver will allow, for example, simple join:
		rebecca.ModelMetadata `tablename:"people JOIN posts ON posts.author_id = people.id"`

		// .. fields are defined here ..
	}

	type PersonCountByAge struct {
		rebecca.ModelMetadata `tablename:"people"`

		// Additionally you can use any expressions as a database mapping for
		// fields, that your driver will allow, for example,
		// count(distinct(field_name)):
		Count int `rebecca:"count(distinct(id))"`
		Age   int `rebecca:"age"`
	}

	// PersonCountByAge is useful, because it can be nicely used with
	// aggregation:
	byAge := []PersonCountByAge{}
	ctx := rebecca.Context{Group: "age"}
	if err := ctx.All(&byAge); err != nil {
		panic(err)
	}

	// At this point byAge contains counts of people per each distinct age.
	// Ordering depends on your chosen driver and database.
	fmt.Print(byAge)
}
