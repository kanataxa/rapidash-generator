# rapidash-generator
Generate automatically [rapidash](https://github.com/knocknote/rapidash) boilerplate functions from sql or go-source

# Features
- Generate `rapidash.Marshaler/Unmarshaler` functions and `rapidash.Struct` from `your go source`

## Not Support(but support future)
- Generate from sql(create tables)

# Install
``` sh
go get github.com/kanataxa/rapidash-generator/cmd/rapi-gen 
```

# Usage
``` sh
rapi-gen -o ${entity}_rapidash.go ${entity}.go
```

``` sh
Usage:
  rapi-gen [OPTIONS]

Application Options:
  -w            force write if file is already exists
  -o, --output= output file name. default: os.Stdout
  -t, --tag=    use tag name (default: db)

Help Options:
  -h, --help    Show this help message
```


For example, you run `rapi-gen` with below input code.

``` go
type User struct {
	ID        uint64     `db:"id"`
	Name      string     `db:"name" json:"name"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
	NonTag    uint
}

type Value uint32
type StrongItem struct {
	ID    uint64  `db:"id"`
	Name  *string `db:"name"`
	Value Value   `db:"value"`
}

type NonDBStruct struct {
	Power uint64
}
```

rapi-gen generates below go code.
``` go
package entity

import (
	"go.knocknote.io/rapidash"
	"golang.org/x/xerrors"
)

func (e *StrongItem) Struct() *rapidash.Struct {
	return rapidash.NewStruct("strong_items").
		FieldUint64("id").
		FieldStringPtr("name").
		FieldUint32("value")
}
func (e *StrongItem) EncodeRapidash(enc rapidash.Encoder) error {
	enc.Uint64("id", e.ID)
	enc.StringPtr("name", e.Name)
	enc.Uint32("value", uint32(e.Value))
	if err := enc.Error(); err != nil {
		return xerrors.Errorf("failed to encode rapidash: %w", err)
	}
	return nil
}
func (e *StrongItem) DecodeRapidash(dec rapidash.Decoder) error {
	e.ID = dec.Uint64("id")
	e.Name = dec.StringPtr("name")
	e.Value = Value(dec.Uint32("value"))
	if err := dec.Error(); err != nil {
		return xerrors.Errorf("failed to decode rapidash: %w", err)
	}
	return nil
}

func (e *User) Struct() *rapidash.Struct {
	return rapidash.NewStruct("users").
		FieldUint64("id").
		FieldString("name").
		FieldTime("created_at").
		FieldTimePtr("updated_at")
}
func (e *User) EncodeRapidash(enc rapidash.Encoder) error {
	enc.Uint64("id", e.ID)
	enc.String("name", e.Name)
	enc.Time("created_at", e.CreatedAt)
	enc.TimePtr("updated_at", e.UpdatedAt)
	if err := enc.Error(); err != nil {
		return xerrors.Errorf("failed to encode rapidash: %w", err)
	}
	return nil
}
func (e *User) DecodeRapidash(dec rapidash.Decoder) error {
	e.ID = dec.Uint64("id")
	e.Name = dec.String("name")
	e.CreatedAt = dec.Time("created_at")
	e.UpdatedAt = dec.TimePtr("updated_at")
	if err := dec.Error(); err != nil {
		return xerrors.Errorf("failed to decode rapidash: %w", err)
	}
	return nil
}
```

# Author
kanataxa(Sota Itoh)

# Notes
This tool is not still stable.
