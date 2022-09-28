package proto

import (
	"bytes"
	"encoding"
	"errors"
	"io"
	"strconv"
	"time"
)

type writer interface {
	io.Writer
	io.ByteWriter
	// WriteString implement io.StringWriter.
	WriteString(s string) (n int, err error)
}

type Writer struct {
	writer
}

func NewWriter(wr writer) *Writer {
	return &Writer{
		writer: wr,
	}
}

func (w *Writer) writeBytes(bs []byte, b *bytes.Buffer) error {
	lbs := strconv.AppendInt(nil, int64(len(bs)), 10)
	if _, err := b.Write(lbs); err != nil {
		return err
	}
	if _, err := b.Write([]byte{EndN}); err != nil {
		return err
	}
	if _, err := b.Write(bs); err != nil {
		return err
	}
	if _, err := b.Write([]byte{EndN}); err != nil {
		return err
	}
	return nil
}

func (w *Writer) WriteArgs(args []interface{}) error {
	var err error
	var buf bytes.Buffer
	for _, arg := range args {
		switch arg := arg.(type) {
		case string:
			err = w.writeBytes([]byte(arg), &buf)
		case []byte:
			err = w.writeBytes(arg, &buf)
		case int:
			bs := strconv.AppendInt(nil, int64(arg), 10)
			err = w.writeBytes(bs, &buf)
		case int8:
			err = w.writeBytes([]byte{byte(arg)}, &buf)
		case int16:
			bs := strconv.AppendInt(nil, int64(arg), 10)
			err = w.writeBytes(bs, &buf)
		case int32:
			bs := strconv.AppendInt(nil, int64(arg), 10)
			err = w.writeBytes(bs, &buf)
		case int64:
			bs := strconv.AppendInt(nil, arg, 10)
			err = w.writeBytes(bs, &buf)
		case uint8:
			err = w.writeBytes([]byte{byte(arg)}, &buf)
		case uint16:
			bs := strconv.AppendUint(nil, uint64(arg), 10)
			err = w.writeBytes(bs, &buf)
		case uint32:
			bs := strconv.AppendUint(nil, uint64(arg), 10)
			err = w.writeBytes(bs, &buf)
		case uint64:
			bs := strconv.AppendUint(nil, arg, 10)
			err = w.writeBytes(bs, &buf)
		case float32:
			bs := strconv.AppendFloat(nil, float64(arg), 'g', -1, 32)
			err = w.writeBytes(bs, &buf)
		case float64:
			bs := strconv.AppendFloat(nil, arg, 'g', -1, 64)
			err = w.writeBytes(bs, &buf)
		case bool:
			if arg {
				err = w.writeBytes([]byte{'1'}, &buf)
			} else {
				err = w.writeBytes([]byte{'0'}, &buf)
			}
		case time.Time:
			bs := strconv.AppendInt(nil, arg.Unix(), 10)
			err = w.writeBytes(bs, &buf)
		case time.Duration:
			bs := strconv.AppendInt(nil, arg.Nanoseconds(), 10)
			err = w.writeBytes(bs, &buf)
		case encoding.BinaryMarshaler:
			b, err := arg.MarshalBinary()
			if err != nil {
				return err
			}
			err = w.writeBytes(b, &buf)
		case nil:
			err = w.writeBytes([]byte{}, &buf)
		default:
			return errors.New("bad arguments type")

		}
		if err != nil {
			return err
		}
	}
	buf.WriteByte(EndN)
	//fmt.Printf("WriteArgs buf: %s", buf.String())
	_, err = w.Write(buf.Bytes())
	return err
}
