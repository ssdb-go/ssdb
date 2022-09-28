package ssdb

import (
	"context"
	"fmt"
	"io"
	"net"
	"testing"
	"time"

	"github.com/ssdb-go/ssdb/internal/proto"
)

var ctx = context.TODO()

type ClientStub struct {
	Cmdable
	resp []byte
}

var initHello = []byte("%1\r\n+proto\r\n:3\r\n")

func NewClientStub(resp []byte) *ClientStub {
	stub := &ClientStub{
		resp: resp,
	}

	stub.Cmdable = NewClient(&Options{
		PoolSize: 128,
		Dialer: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return stub.stubConn(initHello), nil
		},
	})
	return stub
}

func (c *ClientStub) stubConn(init []byte) *ConnStub {
	return &ConnStub{
		init: init,
		resp: c.resp,
	}
}

type ConnStub struct {
	init []byte
	resp []byte
	pos  int
}

func (c *ConnStub) Read(b []byte) (n int, err error) {
	// Return conn.init()
	if len(c.init) > 0 {
		n = copy(b, c.init)
		c.init = c.init[n:]
		return n, nil
	}

	if len(c.resp) == 0 {
		return 0, io.EOF
	}

	if c.pos >= len(c.resp) {
		c.pos = 0
	}
	n = copy(b, c.resp[c.pos:])
	c.pos += n
	return n, nil
}

func (c *ConnStub) Write(b []byte) (n int, err error)  { return len(b), nil }
func (c *ConnStub) Close() error                       { return nil }
func (c *ConnStub) LocalAddr() net.Addr                { return nil }
func (c *ConnStub) RemoteAddr() net.Addr               { return nil }
func (c *ConnStub) SetDeadline(_ time.Time) error      { return nil }
func (c *ConnStub) SetReadDeadline(_ time.Time) error  { return nil }
func (c *ConnStub) SetWriteDeadline(_ time.Time) error { return nil }

type ClientStubFunc func([]byte) *ClientStub

func BenchmarkDecode(b *testing.B) {
	type Benchmark struct {
		name string
		stub ClientStubFunc
	}

	benchmarks := []Benchmark{
		{"server", NewClientStub},
	}

	for _, bench := range benchmarks {
		b.Run(fmt.Sprintf("RespError-%s", bench.name), func(b *testing.B) {
			respError(b, bench.stub)
		})
		b.Run(fmt.Sprintf("RespStatus-%s", bench.name), func(b *testing.B) {
			respStatus(b, bench.stub)
		})
		b.Run(fmt.Sprintf("RespString-%s", bench.name), func(b *testing.B) {
			respString(b, bench.stub)
		})
		b.Run(fmt.Sprintf("RespPipeline-%s", bench.name), func(b *testing.B) {
			respPipeline(b, bench.stub)
		})
		b.Run(fmt.Sprintf("RespTxPipeline-%s", bench.name), func(b *testing.B) {
			respTxPipeline(b, bench.stub)
		})

		// goroutine
		b.Run(fmt.Sprintf("DynamicGoroutine-%s-pool=5", bench.name), func(b *testing.B) {
			dynamicGoroutine(b, bench.stub, 5)
		})
		b.Run(fmt.Sprintf("DynamicGoroutine-%s-pool=20", bench.name), func(b *testing.B) {
			dynamicGoroutine(b, bench.stub, 20)
		})
		b.Run(fmt.Sprintf("DynamicGoroutine-%s-pool=50", bench.name), func(b *testing.B) {
			dynamicGoroutine(b, bench.stub, 50)
		})
		b.Run(fmt.Sprintf("DynamicGoroutine-%s-pool=100", bench.name), func(b *testing.B) {
			dynamicGoroutine(b, bench.stub, 100)
		})

		b.Run(fmt.Sprintf("StaticGoroutine-%s-pool=5", bench.name), func(b *testing.B) {
			staticGoroutine(b, bench.stub, 5)
		})
		b.Run(fmt.Sprintf("StaticGoroutine-%s-pool=20", bench.name), func(b *testing.B) {
			staticGoroutine(b, bench.stub, 20)
		})
		b.Run(fmt.Sprintf("StaticGoroutine-%s-pool=50", bench.name), func(b *testing.B) {
			staticGoroutine(b, bench.stub, 50)
		})
		b.Run(fmt.Sprintf("StaticGoroutine-%s-pool=100", bench.name), func(b *testing.B) {
			staticGoroutine(b, bench.stub, 100)
		})
	}
}

func respError(b *testing.B, stub ClientStubFunc) {
	sdb := stub([]byte("-ERR test error\r\n"))
	respErr := proto.SsdbError("ERR test error")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := sdb.Get(ctx, "key").Err(); err != respErr {
			b.Fatalf("response error, got %q, want %q", err, respErr)
		}
	}
}

func respStatus(b *testing.B, stub ClientStubFunc) {
	sdb := stub([]byte("+OK\r\n"))
	var val interface{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if val = sdb.Set(ctx, "key", "value", 0).Val(); val != "OK" {
			b.Fatalf("response error, got %q, want OK", val)
		}
	}
}

func respString(b *testing.B, stub ClientStubFunc) {
	sdb := stub([]byte("$5\r\nhello\r\n"))
	var val interface{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if val = sdb.Get(ctx, "key").Val(); val != "hello" {
			b.Fatalf("response error, got %q, want hello", val)
		}
	}
}

func respPipeline(b *testing.B, stub ClientStubFunc) {
	sdb := stub([]byte("+OK\r\n$5\r\nhello\r\n:1\r\n"))
	var pipe Pipeliner

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pipe = sdb.Pipeline()
		set := pipe.Set(ctx, "key", "value", 0)
		get := pipe.Get(ctx, "key")
		del := pipe.Del(ctx, "key")
		_, err := pipe.Exec(ctx)
		if err != nil {
			b.Fatalf("response error, got %q, want nil", err)
		}
		if set.Val() != "OK" || get.Val() != "hello" || del.Val() != 1 {
			b.Fatal("response error")
		}
	}
}

func respTxPipeline(b *testing.B, stub ClientStubFunc) {
	sdb := stub([]byte("+OK\r\n+QUEUED\r\n+QUEUED\r\n+QUEUED\r\n*3\r\n+OK\r\n$5\r\nhello\r\n:1\r\n"))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var set *Cmd
		var get *Cmd
		var del *Cmd
		_, err := sdb.TxPipelined(ctx, func(pipe Pipeliner) error {
			set = pipe.Set(ctx, "key", "value", 0)
			get = pipe.Get(ctx, "key")
			del = pipe.Del(ctx, "key")
			return nil
		})
		if err != nil {
			b.Fatalf("response error, got %q, want nil", err)
		}
		if set.Val() != "OK" || get.Val() != "hello" || del.Val() != 1 {
			b.Fatal("response error")
		}
	}
}

func dynamicGoroutine(b *testing.B, stub ClientStubFunc, concurrency int) {
	sdb := stub([]byte("$5\r\nhello\r\n"))
	c := make(chan struct{}, concurrency)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c <- struct{}{}
		go func() {
			if val := sdb.Get(ctx, "key").Val(); val != "hello" {
				panic(fmt.Sprintf("response error, got %q, want hello", val))
			}
			<-c
		}()
	}
	// Here no longer wait for all goroutines to complete, it will not affect the test results.
	close(c)
}

func staticGoroutine(b *testing.B, stub ClientStubFunc, concurrency int) {
	sdb := stub([]byte("$5\r\nhello\r\n"))
	c := make(chan struct{}, concurrency)

	b.ResetTimer()

	for i := 0; i < concurrency; i++ {
		go func() {
			for {
				_, ok := <-c
				if !ok {
					return
				}
				if val := sdb.Get(ctx, "key").Val(); val != "hello" {
					panic(fmt.Sprintf("response error, got %q, want hello", val))
				}
			}
		}()
	}
	for i := 0; i < b.N; i++ {
		c <- struct{}{}
	}
	close(c)
}
