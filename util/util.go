package util

import (
	"bytes"
	"compress/gzip"
	"io"
	"log/slog"
	"net/http"
	"os"
	"peakproxy/common"
	"strings"
)

func ConfigureLogger(logLevel string) (*slog.Logger, *slog.LevelVar) {
	level := new(slog.LevelVar)
	switch strings.ToLower(logLevel) {
	case strings.ToLower(slog.LevelDebug.String()):
		level.Set(slog.LevelDebug)
	case strings.ToLower(slog.LevelInfo.String()):
		level.Set(slog.LevelInfo)
	case strings.ToLower(slog.LevelWarn.String()):
		level.Set(slog.LevelWarn)
	case strings.ToLower(slog.LevelError.String()):
		level.Set(slog.LevelError)
	default:
		common.LogFatalFailedToRecognizeLogLevelErrorWithMessage(logLevel)
	}
	l := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	}))
	slog.SetDefault(l)
	return l, level
}

func CloseWithLog(c io.Closer, l *slog.Logger) {
	if err := c.Close(); err != nil {
		common.LogFailedToCloseError(err, l)
	}
}

func IsGzipEncoded(r *http.Response) bool {
	return r.Header.Get("Content-Encoding") == "gzip"
}

func SetGzipContentEncoding(r *http.Response) {
	r.Header.Set("Content-Encoding", "gzip")
}

func UnsetContentEncoding(r *http.Response) {
	r.Header.Del("Content-Encoding")
}

func ReadAllGzipEncodedWithLog(data []byte, l *slog.Logger) ([]byte, error) {
	if r, err := gzip.NewReader(bytes.NewBuffer(data)); err == nil {
		defer CloseWithLog(r, l)
		return ReadAllWithLog(r, l)
	} else {
		return nil, common.LogAndGetFailedToCreateGzipReaderError(err, l)
	}
}

func WriteAllGzipEncodedWithLog(data []byte, l *slog.Logger) ([]byte, error) {
	b := bytes.Buffer{}
	w := gzip.NewWriter(&b)
	defer CloseWithLog(w, l)
	if _, err := w.Write(data); err != nil {
		return nil, common.LogAndGetFailedToWriteGzipEncodedError(err, l)
	}
	if err := w.Flush(); err == nil {
		return b.Bytes(), nil
	} else {
		return nil, common.LogAndGetFailedToFlushWriterError(err, l)
	}
}

func ReadAllWithLog(r io.Reader, l *slog.Logger) ([]byte, error) {
	if data, err := io.ReadAll(r); err == nil {
		return data, nil
	} else {
		return nil, common.LogAndGetFailedToReadAllError(err, l)
	}
}
