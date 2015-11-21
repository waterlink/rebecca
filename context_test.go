package rebecca_test

import (
	"fmt"
	"reflect"
	"testing"

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

func ExampleContext_All_skipOffsetAlias() {
	type Person struct {
		// ...
	}

	ctx := rebecca.Context{Limit: 20, Offset: 40}
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

func TestContextGetters(t *testing.T) {
	examples := map[string]struct {
		ctx      *rebecca.Context
		action   func(ctx *rebecca.Context) interface{}
		expected interface{}
	}{
		"GetOrder": {
			ctx: &rebecca.Context{Order: "hello DESC"},
			action: func(ctx *rebecca.Context) interface{} {
				return ctx.GetOrder()
			},
			expected: "hello DESC",
		},

		"GetOrder after SetOrder": {
			ctx: &rebecca.Context{Order: "hello DESC"},
			action: func(ctx *rebecca.Context) interface{} {
				x := ctx.SetOrder("id ASC")
				return []string{ctx.GetOrder(), x.GetOrder()}
			},
			expected: []string{"hello DESC", "id ASC"},
		},

		"GetGroup": {
			ctx: &rebecca.Context{Group: "age"},
			action: func(ctx *rebecca.Context) interface{} {
				return ctx.GetGroup()
			},
			expected: "age",
		},

		"GetGroup after SetGroup": {
			ctx: &rebecca.Context{Group: "age"},
			action: func(ctx *rebecca.Context) interface{} {
				x := ctx.SetGroup("level")
				return []string{ctx.GetGroup(), x.GetGroup()}
			},
			expected: []string{"age", "level"},
		},

		"GetLimit": {
			ctx: &rebecca.Context{Limit: 20},
			action: func(ctx *rebecca.Context) interface{} {
				return ctx.GetLimit()
			},
			expected: 20,
		},

		"GetLimit after SetLimit": {
			ctx: &rebecca.Context{Limit: 20},
			action: func(ctx *rebecca.Context) interface{} {
				x := ctx.SetLimit(40)
				return []int{ctx.GetLimit(), x.GetLimit()}
			},
			expected: []int{20, 40},
		},

		"GetSkip from Skip": {
			ctx: &rebecca.Context{Skip: 40},
			action: func(ctx *rebecca.Context) interface{} {
				return ctx.GetSkip()
			},
			expected: 40,
		},

		"GetSkip after SetSkip": {
			ctx: &rebecca.Context{Skip: 40},
			action: func(ctx *rebecca.Context) interface{} {
				x := ctx.SetSkip(20)
				return []int{ctx.GetSkip(), x.GetSkip()}
			},
			expected: []int{40, 20},
		},

		"GetSkip from Offset": {
			ctx: &rebecca.Context{Offset: 40},
			action: func(ctx *rebecca.Context) interface{} {
				return ctx.GetSkip()
			},
			expected: 40,
		},

		"GetSkip ignores Skip if Offset is present": {
			ctx: &rebecca.Context{Skip: 20, Offset: 40},
			action: func(ctx *rebecca.Context) interface{} {
				return ctx.GetSkip()
			},
			expected: 40,
		},

		"GetTx": {
			ctx: &rebecca.Context{},
			action: func(ctx *rebecca.Context) interface{} {
				return ctx.GetTx()
			},
			expected: nil,
		},
	}

	for info, e := range examples {
		t.Log(info)
		actual := e.action(e.ctx)
		if !reflect.DeepEqual(actual, e.expected) {
			t.Errorf("Expected %s to equal %s", actual, e.expected)
		}
	}
}
