package httputil

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"compress/zlib"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// const define
const (
	MaxDecodeBodySize = 8 * 1024 * 1024 // default with 8M limit

	HeaderContentEncoding = "Content-Encoding"

	HeaderContentType = "Content-Type"
)

// NewResponse create a Response
func NewResponse(resp *http.Response) *Response {
	return &Response{
		Response: resp,
	}
}

// Response is response operator
type Response struct {
	*http.Response
	contentEncoding string
}

// ContentType return raw content type from Response
func (r *Response) ContentType() string {
	s := r.Header.Get(HeaderContentType)
	if idx := strings.IndexByte(s, ';'); idx >= 0 {
		return s[:idx]
	}
	return s
}

// Underlying return raw Response
func (r *Response) Underlying() *http.Response {
	return r.Response
}

// EncodeWriteBody write body with encode
func (r *Response) EncodeWriteBody(b []byte, compress bool) error {
	if !compress {
		delete(r.Header, HeaderContentEncoding)
	}
	w := bytes.NewBuffer(nil)
	var wd io.WriteCloser
	var err error
	switch r.contentEncoding = r.Header.Get(HeaderContentEncoding); r.contentEncoding {
	case "gzip":
		wd = gzip.NewWriter(w)
	case "deflate":
		wd, err = flate.NewWriter(w, 1)
	case "br":
		wd = zlib.NewWriter(w)
	case "", "identity":
		wd = NopCloserWriter(w)
	default:
		return fmt.Errorf("unsupport Content-Encoding %s", r.contentEncoding)
	}
	if err != nil {
		return err
	}
	defer wd.Close()

	if _, err = wd.Write(b); err != nil {
		return err
	}

	r.Body = io.NopCloser(w)
	return nil
}

// ReadRawBody read body without decode
func (r *Response) ReadRawBody(limitSize int) ([]byte, error) {
	if r.Body == nil {
		return nil, nil
	}

	b, err := LimitReadAll(r.Body, limitSize, nil)
	r.Body.Close()
	if err != nil {
		return nil, err
	}

	return b, nil
}

// DecodeCloseBody decode body with compression
func (r *Response) DecodeCloseBody(limitSize int) ([]byte, error) {
	if r.Body == nil {
		return nil, nil
	}

	var rd io.Reader
	var err error
	switch r.contentEncoding = r.Header.Get(HeaderContentEncoding); r.contentEncoding {
	case "gzip":
		rd, err = gzip.NewReader(r.Body)
	case "deflate":
		rd = flate.NewReader(r.Body)
	case "br":
		rd, err = zlib.NewReader(r.Body)
	case "", "identity":
		rd = r.Body
	default:
		return nil, fmt.Errorf("unsupport Content-Encoding %s", r.contentEncoding)
	}
	if err != nil {
		return nil, err
	}

	b, err := LimitReadAll(rd, limitSize, nil)
	r.Body.Close()
	if err != nil {
		return nil, err
	}

	return b, nil
}

// LimitReadAll read all data with size limit.
// NOTE: It is extend from io.ReadAll to avoid attack such as gzip bomb.
func LimitReadAll(r io.Reader, limitSize int, buf []byte) ([]byte, error) {
	if limitSize <= 0 {
		limitSize = MaxDecodeBodySize //default size limit
	}
	b := buf[:0]
	for {
		if len(b) >= limitSize {
			return nil, fmt.Errorf("read out of limit size %d/%d", len(b), limitSize)
		}
		if len(b) == cap(b) {
			// Add more capacity (let append pick how much).
			b = append(b, 0)[:len(b)]
		}
		n, err := r.Read(b[len(b):cap(b)])
		b = b[:len(b)+n]
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return b, err
		}
	}
}

// NopCloserWriter is similar with io.NopCloser
func NopCloserWriter(w io.Writer) io.WriteCloser {
	return nopCloserWriter{w: w}
}

type nopCloserWriter struct {
	w io.Writer
}

func (w nopCloserWriter) Write(b []byte) (int, error) {
	return w.w.Write(b)
}

func (w nopCloserWriter) Close() error {
	return nil
}
