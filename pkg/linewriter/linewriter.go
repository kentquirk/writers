package linewriter

// ----- ---- --- -- -
// Copyright 2019, 2020 The Axiom Foundation. All Rights Reserved.
//
// Licensed under the Apache License 2.0 (the "License").  You may not use
// this file except in compliance with the License.  You can obtain a copy
// in the file LICENSE in the source distribution or at
// https://www.apache.org/licenses/LICENSE-2.0.txt
// - -- --- ---- -----


import (
	"bufio"
	"io"
	"unicode/utf8"
)

const newline = 0x0a

// LineWriter wraps an io.Writer and buffers output to it.
//
// It flushes whenever a newline (0x0a, \n) is detected.
//
// The bufio.Writer struct wraps a writer and buffers its
// output. However, it only does this batched write when the
// internal buffer fills. Sometimes, you'd prefer to write
// each line as it's completed, rather than the entire buffer
// at once. Enter LineWriter. It does exactly that.
//
// Like bufio.Writer, a LineWriter's buffer will also be
// flushed when its internal buffer is full. Like
// bufio.Writer, after all data has been written, the
// client should call the Flush method to guarantee that
// all data has been forwarded to the underlying io.Writer.
type LineWriter struct {
	buffer *bufio.Writer
}

// static assert that LineWriter is an io.Writer
var _ io.Writer = (*LineWriter)(nil)

// New creates a new LineWriter
func New(w io.Writer) *LineWriter {
	return &LineWriter{
		buffer: bufio.NewWriter(w),
	}
}

// Write writes the contents of p.
//
// It returns the number of bytes written.
// If n < len(p), it also returns an error explaining
// why the write is short.
func (l *LineWriter) Write(p []byte) (n int, err error) {
	lower := 0

	passthrough := func(upper int, flush bool) error {
		written, err := l.buffer.Write(p[lower:upper])
		n += written
		if err != nil {
			return err
		}

		if flush {
			err = l.buffer.Flush()
			if err != nil {
				return err
			}
		}

		lower = upper
		return nil
	}

	for i, b := range p {
		if b == newline {
			err = passthrough(i+1, true)
			if err != nil {
				return
			}
		}
	}

	if lower < len(p) {
		err = passthrough(len(p), false)
	}
	return
}

// WriteByte writes a single byte
func (l *LineWriter) WriteByte(c byte) error {
	_, err := l.Write([]byte{c})
	return err
}

// WriteRune writes a single Unicode code point.
//
// It returns the number of bytes written and any error.
func (l *LineWriter) WriteRune(r rune) (size int, err error) {
	buf := make([]byte, utf8.UTFMax)
	nbytes := utf8.EncodeRune(buf, r)
	return l.Write(buf[:nbytes])
}

// WriteString writes a string.
//
// It returns the number of bytes written. If the count is
// less than len(s), it also returns an error explaining
// why the write is short.
func (l *LineWriter) WriteString(s string) (int, error) {
	return l.Write([]byte(s))
}

// Flush writes any buffered data to the underlying io.Writer.
func (l *LineWriter) Flush() error {
	return l.buffer.Flush()
}
