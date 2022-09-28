package ssdb_test

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ssdb-go/ssdb"
)

func benchmarkSsdbClient(ctx context.Context, poolSize int) *ssdb.Client {
	client := ssdb.NewClient(&ssdb.Options{
		Addr:         ":8888",
		DialTimeout:  time.Second,
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second,
		PoolSize:     poolSize,
	})
	return client
}

func BenchmarkSsdbPing(b *testing.B) {
	ctx := context.Background()
	sdb := benchmarkSsdbClient(ctx, 10)
	defer sdb.Close()

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if err := sdb.Ping(ctx).Err(); err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkSetGoroutines(b *testing.B) {
	ctx := context.Background()
	sdb := benchmarkSsdbClient(ctx, 10)
	defer sdb.Close()

	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup

		for i := 0; i < 1000; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()

				err := sdb.Set(ctx, "hello", "world", 0).Err()
				if err != nil {
					panic(err)
				}
			}()
		}

		wg.Wait()
	}
}

func BenchmarkSsdbGetNil(b *testing.B) {
	ctx := context.Background()
	client := benchmarkSsdbClient(ctx, 10)
	defer client.Close()

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if err := client.Get(ctx, "key").Err(); err != ssdb.Nil {
				b.Fatal(err)
			}
		}
	})
}

type setStringBenchmark struct {
	poolSize  int
	valueSize int
}

func (bm setStringBenchmark) String() string {
	return fmt.Sprintf("pool=%d value=%d", bm.poolSize, bm.valueSize)
}

func BenchmarkSsdbSetString(b *testing.B) {
	benchmarks := []setStringBenchmark{
		{10, 64},
		{10, 1024},
		{10, 64 * 1024},
		{10, 1024 * 1024},
		{10, 10 * 1024 * 1024},

		{100, 64},
		{100, 1024},
		{100, 64 * 1024},
		{100, 1024 * 1024},
		{100, 10 * 1024 * 1024},
	}
	for _, bm := range benchmarks {
		b.Run(bm.String(), func(b *testing.B) {
			ctx := context.Background()
			client := benchmarkSsdbClient(ctx, bm.poolSize)
			defer client.Close()

			value := strings.Repeat("1", bm.valueSize)

			b.ResetTimer()

			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					err := client.Set(ctx, "key", value, 0).Err()
					if err != nil {
						b.Fatal(err)
					}
				}
			})
		})
	}
}

func BenchmarkSsdbSetGetBytes(b *testing.B) {
	ctx := context.Background()
	client := benchmarkSsdbClient(ctx, 10)
	defer client.Close()

	value := bytes.Repeat([]byte{'1'}, 10000)

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if err := client.Set(ctx, "key", value, 0).Err(); err != nil {
				b.Fatal(err)
			}

			/* got := client.Get(ctx, "key")
			if !bytes.Equal(got, value) {
				b.Fatalf("got != value")
			} */
		}
	})
}

func BenchmarkPipeline(b *testing.B) {
	ctx := context.Background()
	client := benchmarkSsdbClient(ctx, 10)
	defer client.Close()

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := client.Pipelined(ctx, func(pipe ssdb.Pipeliner) error {
				pipe.Set(ctx, "key", "hello", 0)
				//pipe.Expire(ctx, "key", time.Second)
				return nil
			})
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
