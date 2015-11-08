# rebecca

Simple database convenience wrapper for Go language.

Docs are available on http://godoc.org/github.com/waterlink/rebecca

*NOTE* The work is still in progress and not usable.

## Installation

```bash
go get -u github.com/waterlink/rebecca
```

## Usage

```go
import "github.com/waterlink/rebecca"
```

### Designing record

```go
type Person struct {
        rebecca.ModelMetadata `tablename:"people"`

        ID   int    `rebecca:"id" rebecca_primary:"true"`
        Name string `rebecca:"name"`
        Age  int    `rebecca:"age"`
}
```

### Enabling specific driver

```go
import "github.com/waterlink/rebecca/pgdriver"
```

And in code:

```go
rebecca.SetupDriver(pgdriver.NewDriver("postgres://user:pass@host:port/database?sslmode=sslmode"))
```

The same can be done with empty import:

```go
import _ "github.com/waterlink/rebecca/pgdriver/auto"
```

Empty import will fetch connection options from respective environment
variables:

- `REBECCA_PG_USER`
- `REBECCA_PG_PASS`
- `REBECCA_PG_HOST`
- `REBECCA_PG_PORT`
- `REBECCA_PG_DATABASE`
- `REBECCA_PG_SSLMODE`

or, if present: `REBECCA_PG_URL`

To find out each driver's respective environment variables, see its
documentation.

### List of supported drivers

- `github.com/waterlink/rebecca/pgdriver` - driver for postgresql.
- TODO:
  - `github.com/waterlink/rebecca/cassdriver` - driver for cassandra.
  - `github.com/waterlink/rebecca/mysqldriver` - driver for mysql.
  - `github.com/waterlink/rebecca/mongodriver` - driver for mongodb.

### Fetching the record

```go
var p Person
if err := rebecca.Get(ID, &p); err != nil {
        // handle error here
}

// use &p at this point as a found model instance
```

### Saving record

```go
// creates new record
p := &Person{Name: "John Smith", Age: 31}
if err := rebecca.Save(p); err != nil {
        // handle error here
}

// updates the record
p := &Person{}
if err := rebecca.Get(ID, p); err != nil {
        // handle error here
}

p.Age++
if err := rebecca.Save(p); err != nil {
        // handle error here
}
```

### Fetching all records

```go
people := []Person{}
if err := rebecca.All(&people); err != nil {
        // handle error here
}

// people slice will contain found records
```

### Fetching specific records

```go
kids := []Person{}
if err := rebecca.Where(&kids, "age < $1", 12); err != nil {
        // handle error here
}

// kids slice will contain found records
```

### Fetching only first record

```go
kid := &Person{}
if err := rebecca.First(kid, "age < $1", 12); err != nil {
        // handle error here
}

// kid will contain found first record
```

### Removing record

```go
// Given p is *Person:
if err := rebecca.Remove(p); err != nil {
        // handle error here
}
```

### Fetching count for something

First lets define the view for this purpose:

```go
type PeopleCount struct {
          rebecca.ModelMetadata `tablename:"people"`

          Count int `rebecca:"count(id)"`
}
```

And then, lets query for this count:

```go
kidsCount := &PeopleCount{}
if err := rebecca.First(kidsCount, "age < $1", 12); err != nil {
        // handle error here
}

// Now you can use `kidsCount.Count`
```

### Using order, limit and skip

For example, to fetch second 300 records ordered by age.

```go
kidsBatch := []Person{}
ctx := &rebecca.Context{Order: "age ASC", Limit: 300, Skip: 300}
if err := ctx.Where(&kidsBatch, "age < $1", 12); err != nil {
        // handle error here
}

// Now you can use kidsBatch as a second batch of 300 records ordered by age.
```

### Using aggregation

First, lets define our view for aggregation results:

```go
type PeopleByAge struct {
        rebecca.ModelMetadata `tablename:"people"`

        Age   int `rebecca:"age"`
        Count int `rebecca:"count(distinct(id))"`
}
```

Next, lets query this using the `Context`:

```go
peopleByAge := []PeopleByAge{}
ctx := &rebecca.Context{Group: "age"}
if err := ctx.All(&peopleByAge); err != nil {
        // handle error here
}

// Now peopleByAge represents slice of age => count relationship.
```

## Development

After cloning and `cd`-ing into this repo, run `go get ./...` to get all
dependencies going.

Next make sure current codebase is in green state with `go test ./...`.

It is encouraged to use TDD, i.e.: first write/change a test for your change,
then make it green by making necessary modifications to the source.

Make sure your editor runs `goimports` tool on save.

When you are done, make sure you have run whole test suite, `golint ./...` and
`go vet ./...`.

## Contributing

1. Fork it (https://github.com/waterlink/rebecca)
1. Clone it (`git clone git@github.com:my-username/rebecca.git`)
1. `cd` to it (`cd rebecca`)
1. Create your feature branch (`git checkout -b my-new-feature`)
1. Commit your changes (`git commit -am 'Add some new feature'`)
1. Push the branch to your fork (`git push -u origin my-new-feature`)
1. Create a new Pull Request on Github

## Contributors

- [waterlink](https://github.com/waterlink) - Oleksii Fedorov, author,
  maintainer
