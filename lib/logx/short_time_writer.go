package logx

import "io"

type shortTimeWriter struct {
	*shortTimeLogger
	writer io.Writer
}

func NewShortTimeWriter(writer io.Writer, milliseconds int) *shortTimeWriter {
	return &shortTimeWriter{
		shortTimeLogger: newShortTimeLogger(milliseconds),
		writer:          writer,
	}
}

func (w *shortTimeWriter) Write(p []byte) (n int, err error) {
	w.logOrDiscard(func() {
		w.writer.Write(p)
	})
	return len(p), nil
}
