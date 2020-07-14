package slog

import (
	"bytes"
	"fmt"
)

var unknownFile = []byte("???")

var (
	// DefaultFormatter is the default formatter
	DefaultFormatter = ParseFormat("%D %T %l %m")
	// TimedRotatingFormatter is a formatter for TimedRotatingFileWriter
	TimedRotatingFormatter = ParseFormat("[%l %T] %m")
)

// A Formatter containing a sequence of FormatParts.
type Formatter struct {
	formatParts []FormatPart
}

// ParseFormat parses a format string into a formatter.
func ParseFormat(format string) (formatter *Formatter) {
	if format == "" {
		return
	}
	formatter = &Formatter{}
	formatter.findParts([]byte(format))
	formatter.appendByte('\n')
	return
}

/*
Format formats a record to a bytes.Buffer.
Supported format verbs:
	%%: %
	%l: short name of the level
	%T: time string (HH:MM:SS)
	%D: date string (YYYY-mm-DD)
	%s: source code string (filename:line)
	%S: full source code string (/path/filename.go:line)
*/
func (f *Formatter) Format(r *Record, buf *bytes.Buffer) {
	for _, part := range f.formatParts {
		part.Format(r, buf)
	}
}

func (f *Formatter) findParts(format []byte) {
	length := len(format)
	index := bytes.IndexByte(format, '%')
	if index == -1 || index == length-1 {
		if length == 0 {
			return
		}
		if length == 1 {
			f.appendByte(format[0])
		} else {
			f.appendBytes(format)
		}
		return
	}

	if index > 1 {
		f.appendBytes(format[:index])
	} else if index == 1 {
		f.appendByte(format[0])
	}
	switch c := format[index+1]; c {
	case '%':
		f.appendByte('%')
	case 'l':
		f.formatParts = append(f.formatParts, &LevelFormatPart{})
	case 'T':
		f.formatParts = append(f.formatParts, &TimeFormatPart{})
	case 'D':
		f.formatParts = append(f.formatParts, &DateFormatPart{})
	case 'm':
		f.formatParts = append(f.formatParts, &MessageFormatPart{})
	default:
		f.appendBytes([]byte{'%', c})
	}
	f.findParts(format[index+2:])
	return
}

// FormatPart is an interface containing the Format() method.
type FormatPart interface {
	Format(r *Record, buf *bytes.Buffer)
}

// ByteFormatPart is a FormatPart containing a byte.
type ByteFormatPart struct {
	byte byte
}

// Format writes its byte to the buf.
func (p *ByteFormatPart) Format(r *Record, buf *bytes.Buffer) {
	buf.WriteByte(p.byte)
}

// appendByte appends a byte to the formatter.
// If the previous FormatPart is a ByteFormatPart or BytesFormatPart, they will be merged into a BytesFormatPart;
// otherwise a new ByteFormatPart will be created.
func (f *Formatter) appendByte(b byte) {
	parts := f.formatParts
	count := len(parts)
	if count == 0 {
		f.formatParts = append(parts, &ByteFormatPart{byte: b})
	} else {
		var p FormatPart
		lastPart := parts[count-1]
		switch lp := lastPart.(type) {
		case *ByteFormatPart:
			p = &BytesFormatPart{
				bytes: []byte{lp.byte, b},
			}
		case *BytesFormatPart:
			p = &BytesFormatPart{
				bytes: append(lp.bytes, b),
			}
		default:
			p = &ByteFormatPart{byte: b}
		}
		f.formatParts = append(parts, p)
	}
}

// BytesFormatPart is a FormatPart containing a byte slice.
type BytesFormatPart struct {
	bytes []byte
}

// Format writes its bytes to the buf.
func (p *BytesFormatPart) Format(r *Record, buf *bytes.Buffer) {
	buf.Write(p.bytes)
}

// appendBytes appends a byte slice to the formatter.
// If the previous FormatPart is a ByteFormatPart or BytesFormatPart, they will be merged into a BytesFormatPart;
// otherwise a new BytesFormatPart will be created.
func (f *Formatter) appendBytes(bs []byte) {
	parts := f.formatParts
	count := len(parts)
	if count == 0 {
		f.formatParts = append(parts, &BytesFormatPart{bytes: bs})
	} else {
		var p FormatPart
		lastPart := parts[count-1]
		switch lp := lastPart.(type) {
		case *ByteFormatPart:
			p = &BytesFormatPart{
				bytes: append([]byte{lp.byte}, bs...),
			}
		case *BytesFormatPart:
			p = &BytesFormatPart{
				bytes: append(lp.bytes, bs...),
			}
		default:
			p = &BytesFormatPart{bytes: bs}
		}
		f.formatParts = append(parts, p)
	}
}

// LevelFormatPart is a FormatPart of the level placeholder.
type LevelFormatPart struct{}

// Format writes the short level name of the record to the buf.
func (p *LevelFormatPart) Format(r *Record, buf *bytes.Buffer) {
	buf.WriteByte(levelNames[int(r.level)])
}

var (
	//StampNano = "15:04:05.000000000+07:00"
	StampNano = "15:04:05.000000Z0700"
)

// TimeFormatPart is a FormatPart of the time placeholder.
type TimeFormatPart struct{}

// Format writes the time string of the record to the buf.
func (p *TimeFormatPart) Format(r *Record, buf *bytes.Buffer) {
	//hour, min, sec := r.time.Clock()
	//buf.Write(uint2Bytes2(hour))
	//buf.WriteByte(':')
	//buf.Write(uint2Bytes2(min))
	//buf.WriteByte(':')
	//buf.Write(uint2Bytes2(sec))
	buf.WriteString(r.time.Format(StampNano))
}

// DateFormatPart is a FormatPart of the date placeholder.
type DateFormatPart struct{}

// Format writes the date string of the record to the buf.
func (p *DateFormatPart) Format(r *Record, buf *bytes.Buffer) {
	year, mon, day := r.time.Date()
	buf.Write(uint2Bytes4(year))
	buf.WriteByte('-')
	buf.Write(uint2Bytes2(int(mon)))
	buf.WriteByte('-')
	buf.Write(uint2Bytes2(day))
}

// MessageFormatPart is a FormatPart of the message placeholder.
type MessageFormatPart struct{}

// Format writes the formatted message with args to the buf.
func (p *MessageFormatPart) Format(r *Record, buf *bytes.Buffer) {
	if len(r.args) > 0 {
		if r.message == "" {
			fmt.Fprint(buf, r.args...)
		} else {
			fmt.Fprintf(buf, r.message, r.args...)
		}
	} else if r.message != "" {
		buf.WriteString(r.message)
	}
}
