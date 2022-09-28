//go:build gofuzz
// +build gofuzz

package fuzz

import (
	"context"
	"time"

	"github.com/ssdb-go/ssdb"
)

var (
	ctx = context.Background()
	sdb *ssdb.Client
)

func init() {
	sdb = ssdb.NewClient(&ssdb.Options{
		Addr:         ":8888",
		DialTimeout:  10 * time.Second,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		PoolSize:     10,
		PoolTimeout:  10 * time.Second,
	})
}

func Fuzz(data []byte) int {
	arrayLen := len(data)
	if arrayLen < 4 {
		return -1
	}
	maxIter := int(uint(data[0]))
	for i := 0; i < maxIter && i < arrayLen; i++ {
		n := i % arrayLen
		if n == 0 {
			_ = sdb.Set(ctx, string(data[i:]), string(data[i:]), 0).Err()
		} else if n == 1 {
			_, _ = sdb.Get(ctx, string(data[i:])).Result()
		} else if n == 2 {
			_, _ = sdb.Incr(ctx, string(data[i:])).Result()
		} else if n == 3 {
			var cursor uint64
			_, _, _ = sdb.Scan(ctx, cursor, string(data[i:]), 10).Result()
		}
	}
	return 1
}
