package main

import (
	"bytes"
	"encoding/json"
	"github.com/andr31g/peak-detector/peakdetect"
	"io"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path"
	"peakproxy/common"
	"peakproxy/util"
	"reflect"
	"strconv"
)

type PeakProxy struct {
	Proxy *httputil.ReverseProxy
}

type Sample struct {
	UnixTime int64
	Value    string
}

type Samples struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Metric map[string]string `json:"metric"`
			Values []Sample          `json:"values"`
		} `json:"result"`
	} `json:"data"`
}

func GetPeakProxyStructName() string {
	return reflect.TypeOf(PeakProxy{}).String()
}

func NewPeakProxy(targetURI string, iterationCount uint, l *slog.Logger) (*PeakProxy, error) {
	if locator, err := url.Parse(targetURI); err == nil {
		p := &PeakProxy{httputil.NewSingleHostReverseProxy(locator)}
		p.Proxy.ModifyResponse = func(r *http.Response) error {
			return modifyResponse(r, iterationCount, l)
		}
		return p, nil
	} else {
		return nil, common.LogAndGetFailedToParseTargetURIError(err, l)
	}
}

func (s *Sample) MarshalJSON() ([]byte, error) {
	return json.Marshal([]interface{}{s.UnixTime, s.Value})
}

func (s *Sample) UnmarshalJSON(p []byte) error {
	var tmp []json.RawMessage
	if err := json.Unmarshal(p, &tmp); err != nil {
		return err
	}
	if err := json.Unmarshal(tmp[0], &s.UnixTime); err != nil {
		return err
	}
	if err := json.Unmarshal(tmp[1], &s.Value); err != nil {
		return err
	}
	return nil
}

func (p *PeakProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.Proxy.ServeHTTP(w, r)
}

func modifyResponse(r *http.Response, iterationCount uint, l *slog.Logger) error {
	l.Debug("modifyResponse", "request-url", r.Request.URL.String())

	endpoint := path.Base(r.Request.URL.Path)
	if endpoint != "query_range" {
		l.Debug("only query_range endpoints are supported; returning")
		return nil
	}

	var data []byte
	if body, err := util.ReadAllWithLog(r.Body, l); err == nil {
		data = body
	} else {
		return common.LogAndGetFailedToReadResponseBodyError(err, l)
	}

	if util.IsGzipEncoded(r) {
		if uncompressed, err := util.ReadAllGzipEncodedWithLog(data, l); err == nil {
			data = uncompressed
		} else {
			return common.LogAndGetFailedToReadGzipEncodedError(err, l)
		}
	}

	if body, err := handleResponseJSONAndDetectPeaks(data, iterationCount, l); err == nil {
		return modifyResponseBody(r, body, util.IsGzipEncoded(r), l)
	} else {
		return common.LogAndGetFailedToHandleResponseError(err, l)
	}
}

func modifyResponseBody(r *http.Response, data []byte, gzipEncode bool, l *slog.Logger) error {
	var body []byte
	if gzipEncode {
		if encoded, err := util.WriteAllGzipEncodedWithLog(data, l); err == nil {
			util.SetGzipContentEncoding(r)
			body = encoded
		} else {
			return common.LogAndGetFailedToModifyResponseBodyError(err, l)
		}
	} else {
		util.UnsetContentEncoding(r)
		body = data
	}
	r.Body = io.NopCloser(bytes.NewBuffer(body))
	r.ContentLength = int64(len(body))
	r.Header.Set("Content-Length", strconv.Itoa(len(body)))
	return nil
}

func handleResponseJSONAndDetectPeaks(data []byte, iterationCount uint, l *slog.Logger) ([]byte, error) {
	var s Samples
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, common.LogAndGetFailedToUnmarshalJSONError(err, l)
	}
	if err := detectPeaks(&s, iterationCount, l); err != nil {
		return nil, common.LogAndGetFailedToRunPeakDetectError(err, l)
	}
	if encoded, err := json.Marshal(&s); err == nil {
		return encoded, nil
	} else {
		return nil, common.LogAndGetFailedToMarshalJSONError(err, l)
	}
}

func detectPeaks(s *Samples, iterationCount uint, l *slog.Logger) error {
	if iterationCount == 0 {
		return common.LogAndGetFailedToRunPeakDetectErrorWithMessage("zero iterations specified", l)
	}
	for _, result := range s.Data.Result {
		samples := make([]float64, len(result.Values))
		for i, sample := range result.Values {
			if f, err := strconv.ParseFloat(sample.Value, 64); err == nil {
				samples[i] = f
			} else {
				return common.LogAndGetFailedToConvertSampleToFloatError(err, l)
			}
		}
		var peaks []float64
		if iterationCount == 1 {
			primary := peakdetect.DetectPeaks(samples)
			peaks = primary.InflateWithCount(len(samples), &primary)
		} else {
			if secondary, ok := peakdetect.IteratePeakDetect(iterationCount, samples); ok {
				peaks = secondary.InflateWithCount(len(samples), peakdetect.PrimaryValuesOnly[float64](&secondary))
			}
		}
		for i, _ := range result.Values {
			peak := strconv.FormatFloat(peaks[i], 'f', -1, 64)
			result.Values[i].Value = peak
		}
	}
	return nil
}
