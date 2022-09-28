package ssdb

import (
	"context"
	"time"

	"github.com/ssdb-go/ssdb/internal"
)

// KeepTTL is a ssdb KEEPTTL option to keep existing TTL, it requires your ssdb-server version >= 6.0,
// otherwise you will receive an error: (error) ERR syntax error.
// For example:
//
//    sdb.Set(ctx, key, value, ssdb.KeepTTL)
const KeepTTL = -1

func usePrecise(dur time.Duration) bool {
	return dur < time.Second || dur%time.Second != 0
}

func formatMs(ctx context.Context, dur time.Duration) int64 {
	if dur > 0 && dur < time.Millisecond {
		internal.Logger.Printf(
			ctx,
			"specified duration is %s, but minimal supported value is %s - truncating to 1ms",
			dur, time.Millisecond,
		)
		return 1
	}
	return int64(dur / time.Millisecond)
}

func formatSec(ctx context.Context, dur time.Duration) int64 {
	if dur > 0 && dur < time.Second {
		internal.Logger.Printf(
			ctx,
			"specified duration is %s, but minimal supported value is %s - truncating to 1s",
			dur, time.Second,
		)
		return 1
	}
	return int64(dur / time.Second)
}

func appendArgs(dst, src []interface{}) []interface{} {
	if len(src) == 1 {
		return appendArg(dst, src[0])
	}

	dst = append(dst, src...)
	return dst
}

func appendArg(dst []interface{}, arg interface{}) []interface{} {
	switch arg := arg.(type) {
	case []string:
		for _, s := range arg {
			dst = append(dst, s)
		}
		return dst
	case []interface{}:
		dst = append(dst, arg...)
		return dst
	case map[string]interface{}:
		for k, v := range arg {
			dst = append(dst, k, v)
		}
		return dst
	case map[string]string:
		for k, v := range arg {
			dst = append(dst, k, v)
		}
		return dst
	default:
		return append(dst, arg)
	}
}

type Cmdable interface {
	Pipeline() Pipeliner
	Pipelined(ctx context.Context, fn func(Pipeliner) error) ([]Cmder, error)

	TxPipelined(ctx context.Context, fn func(Pipeliner) error) ([]Cmder, error)
	TxPipeline() Pipeliner

	// db server
	DBSize(ctx context.Context) *Cmd
	DBInfo(ctx context.Context) *Cmd
	Ping(ctx context.Context) *Cmd

	// key-value
	Set(ctx context.Context, key string, val interface{}, ttl ...int64) *Cmd
	SetNX(ctx context.Context, key string, val interface{}) *Cmd
	Get(ctx context.Context, key string) *Cmd
	GetSet(ctx context.Context, key string, val interface{}) *Cmd
	Del(ctx context.Context, key string) *Cmd
	HSet(ctx context.Context, key string, val ...interface{}) *Cmd
	Expire(ctx context.Context, key string, ttl int64) *Cmd
	Exists(ctx context.Context, key string) *Cmd
	TTL(ctx context.Context, key string) *Cmd
	Incr(ctx context.Context, key string, num int64) *Cmd
	/* MultiSet(ctx context.Context, kvs map[string]interface{}) *Cmd
	MultiGet(ctx context.Context, key ...string) *Cmd
	MultiGetSlice(ctx context.Context, key ...string) *Cmd
	MultiGetArray(ctx context.Context, key []string) *Cmd
	MultiGetSliceArray(ctx context.Context, key []string) *Cmd
	MultiDel(ctx context.Context, key ...string) *Cmd
	Setbit(ctx context.Context, key string, offset int64, bit int) *Cmd
	Getbit(ctx context.Context, key string, offset int64) *Cmd
	BitCount(ctx context.Context, key string, start int64, end int64) *Cmd
	CountBit(ctx context.Context, key string, start int64, size int64) *Cmd
	Substr(ctx context.Context, key string, start int64, size ...int64) *Cmd
	StrLen(ctx context.Context, key string) *Cmd
	Keys(ctx context.Context, keyStart, keyEnd string, limit int64) *Cmd
	RKeys(ctx context.Context, keyStart, keyEnd string, limit int64) *Cmd
	Scan(ctx context.Context, keyStart, keyEnd string, limit int64) *Cmd
	RScan(ctx context.Context, keyStart, keyEnd string, limit int64) */
}

type StatefulCmdable interface {
	Cmdable
	Auth(ctx context.Context, password string) *Cmd
	AuthACL(ctx context.Context, username, password string) *Cmd
	ClientSetName(ctx context.Context, name string) *Cmd
}

var (
	_ Cmdable = (*Client)(nil)
)

type cmdable func(ctx context.Context, cmd Cmder) error

type statefulCmdable func(ctx context.Context, cmd Cmder) error

//------------------------------------------------------------------------------

func (c statefulCmdable) Auth(ctx context.Context, password string) *Cmd {
	cmd := NewCmd(ctx, "auth", password)
	_ = c(ctx, cmd)
	return cmd
}

// AuthACL Perform an AUTH command, using the given user and pass.
// Should be used to authenticate the current connection with one of the connections defined in the ACL list
// when connecting to a ssdb 6.0 instance, or greater, that is using the ssdb ACL system.
func (c statefulCmdable) AuthACL(ctx context.Context, username, password string) *Cmd {
	cmd := NewCmd(ctx, "auth", password)
	_ = c(ctx, cmd)
	return cmd
}

// ClientSetName assigns a name to the connection.
func (c statefulCmdable) ClientSetName(ctx context.Context, name string) *Cmd {
	cmd := NewCmd(ctx, "client", "setname", name)
	_ = c(ctx, cmd)
	return cmd
}

//------------------------------------------------------------------------------
func (c cmdable) DBSize(ctx context.Context) *Cmd {
	cmd := NewCmd(ctx, "dbsize")
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) DBInfo(ctx context.Context) *Cmd {
	cmd := NewCmd(ctx, "info")
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) Set(ctx context.Context, key string, val interface{}, ttl ...int64) *Cmd {
	cmd := NewCmd(ctx, "command")
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) Ping(ctx context.Context) *Cmd {
	cmd := NewCmd(ctx, "version")
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) SetNX(ctx context.Context, key string, val interface{}) *Cmd {
	cmd := NewCmd(ctx, "command")
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) Get(ctx context.Context, key string) *Cmd {
	cmd := NewCmd(ctx, "command")
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) GetSet(ctx context.Context, key string, val interface{}) *Cmd {
	cmd := NewCmd(ctx, "command")
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) Del(ctx context.Context, key string) *Cmd {
	cmd := NewCmd(ctx, "command")
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) HSet(ctx context.Context, key string, val ...interface{}) *Cmd {
	cmd := NewCmd(ctx, "command")
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) Scan(ctx context.Context, keyStart, keyEnd string, limit int64) *Cmd {
	cmd := NewCmd(ctx, "command")
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) Expire(ctx context.Context, key string, ttl int64) *Cmd {
	cmd := NewCmd(ctx, "command")
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) Exists(ctx context.Context, key string) *Cmd {
	cmd := NewCmd(ctx, "command")
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) TTL(ctx context.Context, key string) *Cmd {
	cmd := NewCmd(ctx, "command")
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) Incr(ctx context.Context, key string, num int64) *Cmd {
	cmd := NewCmd(ctx, "command")
	_ = c(ctx, cmd)
	return cmd
}
