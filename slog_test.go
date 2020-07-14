package slog

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

type memoryWriter struct {
	bytes.Buffer
}

func (w *memoryWriter) Close() error {
	w.Buffer.Reset()
	return nil
}

func TestSample(t *testing.T) {
	w, err := NewFileWriter("test.txt")
	//golog.NewTimedRotatingFileWriter()
	if err != nil {
		t.Error(err)
		return
	}
	h := NewHandler(InfoLevel, DefaultFormatter)
	h.AddWriter(w)
	l := NewLogger(InfoLevel)
	l.AddHandler(h)
	SetDefaultLogger(l)

	for i := 0; i < 10; i++ {
		Infof("abc %v", i)
	}
}

func TestLogFuncs(t *testing.T) {
	w := &memoryWriter{}
	h := NewHandler(InfoLevel, DefaultFormatter)
	h.AddWriter(w)
	l := NewLogger(InfoLevel)
	l.AddHandler(h)
	SetDefaultLogger(l)

	l.Debug("test")
	if w.Buffer.Len() != 0 {
		t.Error("memoryWriter is not empty")
	}
	l.Debugf("test")
	if w.Buffer.Len() != 0 {
		t.Error("memoryWriter is not empty")
	}

	l.Info("test")
	if w.Buffer.Len() == 0 {
		t.Error("memoryWriter is empty")
	}
	w.Buffer.Reset()

	l.Infof("test")
	if w.Buffer.Len() == 0 {
		t.Error("memoryWriter is empty")
	}
	w.Buffer.Reset()

	l.Error("test")
	if w.Buffer.Len() == 0 {
		t.Error("memoryWriter is empty")
	}
	w.Buffer.Reset()

	l.Errorf("test")
	if w.Buffer.Len() == 0 {
		t.Error("memoryWriter is empty")
	}
	l.Close()

	h = NewHandler(ErrorLevel, DefaultFormatter)
	h.AddWriter(w)
	l = NewLogger(ErrorLevel)
	l.AddHandler(h)
	SetDefaultLogger(l)

	l.Info("test")
	if w.Buffer.Len() != 0 {
		t.Error("memoryWriter is not empty")
	}
	w.Buffer.Reset()

	l.Error("test")
	if w.Buffer.Len() == 0 {
		t.Error("memoryWriter is empty")
	}
	l.Close()
}

func BenchmarkBufferedFileLogger(b *testing.B) {
	path := filepath.Join(os.TempDir(), "test.log")
	os.Remove(path)
	w, err := NewBufferedFileWriter(path)
	if err != nil {
		b.Error(err)
	}
	h := NewHandler(InfoLevel, DefaultFormatter)
	h.AddWriter(w)
	l := NewLogger(InfoLevel)
	l.AddHandler(h)
	SetDefaultLogger(l)

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			Infof("test")
		}
	})

	b.StopTimer()

	l.Close()
}

func BenchmarkDiscardLogger(b *testing.B) {
	w := NewDiscardWriter()
	h := NewHandler(InfoLevel, DefaultFormatter)
	h.AddWriter(w)
	l := NewLogger(InfoLevel)
	l.AddHandler(h)
	SetDefaultLogger(l)

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			Infof("test")
		}
	})
	l.Close()
}

func BenchmarkNopLog(b *testing.B) {
	w := NewDiscardWriter()
	h := NewHandler(InfoLevel, DefaultFormatter)
	h.AddWriter(w)
	l := NewLogger(InfoLevel)
	l.AddHandler(h)
	SetDefaultLogger(l)

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			Debugf("test")
		}
	})
	l.Close()
}

func BenchmarkMultiLevels(b *testing.B) {
	w := NewDiscardWriter()
	dh := NewHandler(DebugLevel, DefaultFormatter)
	dh.AddWriter(w)
	ih := NewHandler(InfoLevel, DefaultFormatter)
	ih.AddWriter(w)
	wh := NewHandler(WarnLevel, DefaultFormatter)
	wh.AddWriter(w)
	eh := NewHandler(ErrorLevel, DefaultFormatter)
	eh.AddWriter(w)
	ch := NewHandler(CritLevel, DefaultFormatter)
	ch.AddWriter(w)

	l := NewLogger(WarnLevel)
	l.AddHandler(dh)
	l.AddHandler(ih)
	l.AddHandler(wh)
	l.AddHandler(eh)
	l.AddHandler(ch)
	SetDefaultLogger(l)

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			Debugf("test")
			Infof("test")
			Warnf("test")
			Errorf("test")
			Critf("test")
		}
	})
	l.Close()
}
