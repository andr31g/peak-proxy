package common

import (
	errs "errors"
	"fmt"
	"log"
	"log/slog"
	"os"
)

var (
	wrapperErrorStr                 = "wrapper error"
	wrappedErrorStr                 = "wrapped error"
	errorOccurredStr                = "error occurred"
	FailedToClose                   = errs.New("failed to close Closeable")
	FailedToReadAllErr              = errs.New("failed to read all")
	FailedToFlushWriterErr          = errs.New("failed to flush Writer")
	FailedToRunPeakDetectErr        = errs.New("failed to run peak detect")
	FailedToMarshalJSONErr          = errs.New("failed to marshal JSON")
	FailedToUnmarshalJSONErr        = errs.New("failed to unmarshal JSON")
	FailedToHandleResponseErr       = errs.New("failed to handle HTTP response")
	FailedToModifyResponseBodyErr   = errs.New("failed to modify HTTP response body")
	FailedToReadResponseBodyErr     = errs.New("failed to read HTTP response body")
	FailedToCreateGzipReaderErr     = errs.New("failed to create gzip reader")
	FailedToReadGzipEncodedErr      = errs.New("failed to read gzip-encoded content")
	FailedToWriteGzipEncodedErr     = errs.New("failed to write gzip-encoded content")
	FailedToConvertSampleToFloatErr = errs.New("failed to convert sample to float")
	FailedToParseTargetURIErr       = errs.New("failed to parse the target URI")
	FailedToCreateInstanceErr       = errs.New("failed to create instance")
	FailedToStartHTTPServerErr      = errs.New("failed to start HTTP server")
	FailedToRecognizeLogLevelErr    = errs.New("failed to recognize log level")
)

func LogFailedToCloseError(err error, l *slog.Logger) {
	logWrappedError(FailedToClose, err, l)
}

func LogAndGetFailedToRunPeakDetectError(err error, l *slog.Logger) error {
	return LogAndGetWrappedError(FailedToRunPeakDetectErr, err, l)
}

func LogAndGetFailedToRunPeakDetectErrorWithMessage(message string, l *slog.Logger) error {
	return LogAndGetWrappedError(FailedToRunPeakDetectErr, errs.New(message), l)
}

func LogAndGetFailedToReadAllError(err error, l *slog.Logger) error {
	return LogAndGetWrappedError(FailedToReadAllErr, err, l)
}

func LogAndGetFailedToFlushWriterError(err error, l *slog.Logger) error {
	return LogAndGetWrappedError(FailedToFlushWriterErr, err, l)
}

func LogAndGetFailedToHandleResponseError(err error, l *slog.Logger) error {
	return LogAndGetWrappedError(FailedToHandleResponseErr, err, l)
}

func LogAndGetFailedToReadResponseBodyError(err error, l *slog.Logger) error {
	return LogAndGetWrappedError(FailedToReadResponseBodyErr, err, l)
}

func LogAndGetFailedToModifyResponseBodyError(err error, l *slog.Logger) error {
	return LogAndGetWrappedError(FailedToModifyResponseBodyErr, err, l)
}

func LogAndGetFailedToCreateGzipReaderError(err error, l *slog.Logger) error {
	return LogAndGetWrappedError(FailedToCreateGzipReaderErr, err, l)
}

func LogAndGetFailedToReadGzipEncodedError(err error, l *slog.Logger) error {
	return LogAndGetWrappedError(FailedToReadGzipEncodedErr, err, l)
}

func LogAndGetFailedToWriteGzipEncodedError(err error, l *slog.Logger) error {
	return LogAndGetWrappedError(FailedToWriteGzipEncodedErr, err, l)
}

func LogAndGetFailedToMarshalJSONError(err error, l *slog.Logger) error {
	return LogAndGetWrappedError(FailedToMarshalJSONErr, err, l)
}

func LogAndGetFailedToUnmarshalJSONError(err error, l *slog.Logger) error {
	return LogAndGetWrappedError(FailedToUnmarshalJSONErr, err, l)
}

func LogAndGetFailedToConvertSampleToFloatError(err error, l *slog.Logger) error {
	return LogAndGetWrappedError(FailedToConvertSampleToFloatErr, err, l)
}

func LogAndGetFailedToParseTargetURIError(err error, l *slog.Logger) error {
	return LogAndGetWrappedError(FailedToParseTargetURIErr, err, l)
}

func LogFatalFailedToRecognizeLogLevelErrorWithMessage(message string) {
	log.Fatalf(GetWrappedError(FailedToRecognizeLogLevelErr, errs.New(message)).Error())
}

func LogFatalFailedToStartHTTPServerError(err error, l *slog.Logger) {
	logWrappedError(FailedToStartHTTPServerErr, err, l)
	os.Exit(1)
}

func LogFatalFailedToCreateInstanceErrorWithMessage(err error, message string, l *slog.Logger) {
	logWrappedError(FailedToCreateInstanceErr, GetWrappedError(errs.New(message), err), l)
	os.Exit(1)
}

func LogAndGetWrappedError(wrapper error, wrapped error, l *slog.Logger) error {
	logWrappedError(wrapper, wrapped, l)
	return GetWrappedError(wrapper, wrapped)
}

func GetWrappedError(wrapper error, wrapped error) error {
	return fmt.Errorf(wrapperErrorStr+"=%w\n"+wrappedErrorStr+"=%w", wrapper, wrapped)
}

func logWrappedError(wrapper error, wrapped error, l *slog.Logger) {
	l.Error(errorOccurredStr, wrapperErrorStr, wrapper.Error(), wrappedErrorStr, wrapped.Error())
}
