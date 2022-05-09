package httputil_test

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/quanxiang-cloud/cabin/lib/httputil"
)

var gzipBombData = testGenGzipBomb()

const testShowDebug = false

func testGenGzipBomb() []byte {
	w := bytes.NewBuffer(nil)
	wd := gzip.NewWriter(w)
	var data [1024]byte
	for i := 0; i < httputil.MaxDecodeBodySize/len(data)+1; i++ {
		wd.Write(data[:])
	}
	wd.Close()
	return w.Bytes()
}

func TestGzipBomb(t *testing.T) {
	resp := &http.Response{
		Body:   io.NopCloser(bytes.NewReader(gzipBombData)),
		Header: http.Header{},
	}
	resp.Header.Set(httputil.HeaderContentEncoding, "gzip")
	rd := httputil.NewResponse(resp)
	b, err := rd.DecodeCloseBody(0)
	if err == nil {
		t.Fatalf("expect error but got nil")
	}
	if testShowDebug {
		fmt.Println("decode bomb", len(gzipBombData), "=>", len(b), err)
	}
}

func TestModifyBody(t *testing.T) {
	resp := &http.Response{
		Body:   io.NopCloser(bytes.NewReader(gzipBombData)),
		Header: http.Header{},
	}
	resp.Header.Set(httputil.HeaderContentEncoding, "gzip")
	rd := httputil.NewResponse(resp)

	origin := `{"a":"1","b":"2"}`
	err := rd.EncodeWriteBody([]byte(origin), true)
	b, err := rd.DecodeCloseBody(0)
	fmt.Println("origin", origin)
	fmt.Println("decoded", string(b), err)
	var d map[string]interface{}
	err = json.Unmarshal(b, &d)

	delete(d, "a")
	b, err = json.Marshal(d)
	err = rd.EncodeWriteBody(b, true)
	b, err = rd.DecodeCloseBody(0)
	fmt.Println("remove a", string(b), err)

	rd.EncodeWriteBody(b, true)
	bb, err := rd.ReadRawBody(0)
	fmt.Println("read encoded raw", string(bb), err)
	rd.EncodeWriteBody(b, false)
	bb, err = rd.ReadRawBody(0)
	fmt.Println("read raw", string(bb), err)
}
