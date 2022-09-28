package main

import (
	"context"

	"github.com/ssdb-go/ssdb"
)

type Model struct {
	Str1    string   `ssdb:"str1"`
	Str2    string   `ssdb:"str2"`
	Int     int      `ssdb:"int"`
	Bool    bool     `ssdb:"bool"`
	Ignored struct{} `ssdb:"-"`
}

func main() {
	ctx := context.Background()

	sdb := ssdb.NewClient(&ssdb.Options{
		Addr: ":8888",
	})

	// Set some fields.
	if _, err := sdb.Pipelined(ctx, func(sdb ssdb.Pipeliner) error {
		sdb.HSet(ctx, "key", "str1", "hello")
		sdb.HSet(ctx, "key", "str2", "world")
		sdb.HSet(ctx, "key", "int", 123)
		sdb.HSet(ctx, "key", "bool", 1)
		return nil
	}); err != nil {
		panic(err)
	}

	/* var model1, model2 Model

	// Scan all fields into the model.
	if err := sdb.HGetAll(ctx, "key").Scan(&model1); err != nil {
		panic(err)
	}

	// Or scan a subset of the fields.
	if err := sdb.HMGet(ctx, "key", "str1", "int").Scan(&model2); err != nil {
		panic(err)
	}

	spew.Dump(model1)
	spew.Dump(model2) */
}
