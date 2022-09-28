package ssdb

import (
	"context"
	"io"
	"net"
	"strings"

	"github.com/ssdb-go/ssdb/internal/pool"
	"github.com/ssdb-go/ssdb/internal/proto"
)

// ErrClosed performs any operation on the closed client will return this error.
var ErrClosed = pool.ErrClosed

type Error interface {
	error

	// SsdbError is a no-op function but
	// serves to distinguish types that are Ssdb
	// errors from ordinary errors: a type is a
	// Ssdb error if it has a SsdbError method.
	SsdbError()
}

var _ Error = proto.SsdbError("")

func shouldRetry(err error, retryTimeout bool) bool {
	switch err {
	case io.EOF, io.ErrUnexpectedEOF:
		return true
	case nil, context.Canceled, context.DeadlineExceeded:
		return false
	}

	if v, ok := err.(timeoutError); ok {
		if v.Timeout() {
			return retryTimeout
		}
		return true
	}

	s := err.Error()
	if s == "ERR max number of clients reached" {
		return true
	}
	if strings.HasPrefix(s, "LOADING ") {
		return true
	}
	if strings.HasPrefix(s, "READONLY ") {
		return true
	}
	if strings.HasPrefix(s, "CLUSTERDOWN ") {
		return true
	}
	if strings.HasPrefix(s, "TRYAGAIN ") {
		return true
	}

	return false
}

func isSsdbError(err error) bool {
	_, ok := err.(proto.SsdbError)
	return ok
}

func isBadConn(err error, allowTimeout bool, addr string) bool {
	switch err {
	case nil:
		return false
	case context.Canceled, context.DeadlineExceeded:
		return true
	}

	if isSsdbError(err) {
		switch {
		case isReadOnlyError(err):
			// Close connections in read only state in case domain addr is used
			// and domain resolves to a different Ssdb Server. See #790.
			return true
		case isMovedSameConnAddr(err, addr):
			// Close connections when we are asked to move to the same addr
			// of the connection. Force a DNS resolution when all connections
			// of the pool are recycled
			return true
		default:
			return false
		}
	}

	if allowTimeout {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			return false
		}
	}

	return true
}

func isMovedError(err error) (moved bool, ask bool, addr string) {
	if !isSsdbError(err) {
		return
	}

	s := err.Error()
	switch {
	case strings.HasPrefix(s, "MOVED "):
		moved = true
	case strings.HasPrefix(s, "ASK "):
		ask = true
	default:
		return
	}

	ind := strings.LastIndex(s, " ")
	if ind == -1 {
		return false, false, ""
	}
	addr = s[ind+1:]
	return
}

func isLoadingError(err error) bool {
	return strings.HasPrefix(err.Error(), "LOADING ")
}

func isReadOnlyError(err error) bool {
	return strings.HasPrefix(err.Error(), "READONLY ")
}

func isMovedSameConnAddr(err error, addr string) bool {
	ssdbError := err.Error()
	if !strings.HasPrefix(ssdbError, "MOVED ") {
		return false
	}
	return strings.HasSuffix(ssdbError, " "+addr)
}

//------------------------------------------------------------------------------

type timeoutError interface {
	Timeout() bool
}
