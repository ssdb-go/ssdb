package proto

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
)

// ssdb resp protocol data type.
const (
	EndN  = '\n'   // +<string>\n
	EndNN = "\n\n" // +<string>\n\n
	EndOK = "ok"   // +<string>ok
)

// Not used temporarily.
// Ssdb has not used these two data types for the time being, and will implement them later.
// Streamed           = "EOF:"
// StreamedAggregated = '?'

//------------------------------------------------------------------------------

const Nil = SsdbError("Ssdb: nil") // nolint:errname

type SsdbError string

func (e SsdbError) Error() string { return string(e) }

func (SsdbError) SsdbError() {}

func ParseErrorReply(line []byte) error {
	return SsdbError(line[1:])
}

//------------------------------------------------------------------------------

type Reader struct {
	rd *bufio.Reader
}

func NewReader(rd io.Reader) *Reader {
	return &Reader{
		rd: bufio.NewReader(rd),
	}
}

func (r *Reader) Buffered() int {
	return r.rd.Buffered()
}

func (r *Reader) Peek(n int) ([]byte, error) {
	return r.rd.Peek(n)
}

func (r *Reader) Reset(rd io.Reader) {
	fmt.Println("golf: Reset")
	r.rd.Reset(rd)
}

// readLine returns an error if:
//   - there is a pending read error.
func (r *Reader) readLine() ([]byte, error) {
	b, err := r.rd.ReadSlice('\n')
	if err != nil {
		if err != bufio.ErrBufferFull {
			return nil, err
		}

		full := make([]byte, len(b))
		copy(full, b)

		b, err = r.rd.ReadBytes('\n')
		if err != nil {
			return nil, err
		}

		full = append(full, b...) //nolint:makezero
		b = full
	}
	if len(b) < 1 || b[len(b)-1] != '\n' {
		return nil, fmt.Errorf("ssdb: invalid reply: %q", b)
	}
	return b, nil
}

// readLines returns an error if:
//   - there is a pending read line error.
func (r *Reader) readLines() ([]byte, error) {
	var data []byte
	var lastBytes []byte
	for {
		b, err := r.readLine()
		//fmt.Println("readLine data: ", string(b), " readLine err:", err)
		if err != nil {
			return nil, err
		}

		if string(b) == "\n" && string(lastBytes) != "1" {
			//fmt.Println("readLines end")
			break
		}

		data = append(data, b...) //nolint:makezero
		lastBytes = b
	}

	return data, nil
}

func (r *Reader) ReadReply() (interface{}, error) {
	buf, err := r.readLines()
	if err != nil {
		return nil, err
	}

	resp := []string{}
	bufArray := bytes.Split(buf, []byte("\n"))
	for i := 1; i < len(bufArray); i = i + 2 {
		resp = append(resp, string(bufArray[i]))
	}

	return resp, nil
}
