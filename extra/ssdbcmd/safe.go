//go:build appengine
// +build appengine

package ssdbcmd

func String(b []byte) string {
	return string(b)
}

func Bytes(s string) []byte {
	return []byte(s)
}
