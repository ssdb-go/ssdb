package ssdb_test

import (
	"context"
	"fmt"
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
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		PoolSize:     10,
		PoolTimeout:  30 * time.Second,
	})
}

func ExampleNewClient() {
	sdb := ssdb.NewClient(&ssdb.Options{
		Addr:     "localhost:8888", // use default Addr
		Password: "",               // no password set
		DB:       0,                // use default DB
	})

	pong, err := sdb.Ping(ctx).Result()
	fmt.Println(pong, err)
	// Output: PONG <nil>
}

func ExampleParseURL() {
	opt, err := ssdb.ParseURL("ssdb://:qwerty@localhost:8888/1?dial_timeout=5s")
	if err != nil {
		panic(err)
	}
	fmt.Println("addr is", opt.Addr)
	fmt.Println("db is", opt.DB)
	fmt.Println("password is", opt.Password)
	fmt.Println("dial timeout is", opt.DialTimeout)

	// Create client as usually.
	_ = ssdb.NewClient(opt)

	// Output: addr is localhost:8888
	// db is 1
	// password is qwerty
	// dial timeout is 5s
}

func ExampleClient() {
	err := sdb.Set(ctx, "key", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := sdb.Get(ctx, "key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key", val)

	val2, err := sdb.Get(ctx, "missing_key").Result()
	if err == ssdb.Nil {
		fmt.Println("missing_key does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("missing_key", val2)
	}
	// Output: key value
	// missing_key does not exist
}

func ExampleConn() {
	conn := sdb.Conn()

	err := conn.ClientSetName(ctx, "foobar").Err()
	if err != nil {
		panic(err)
	}

	// Open other connections.
	for i := 0; i < 10; i++ {
		go sdb.Ping(ctx)
	}
	// Output: foobar
}

func ExampleClient_Set() {
	// Last argument is expiration. Zero means the key has no
	// expiration time.
	err := sdb.Set(ctx, "key", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	// key2 will expire in an hour.
	err = sdb.Set(ctx, "key2", "value", 0).Err()
	if err != nil {
		panic(err)
	}
}
