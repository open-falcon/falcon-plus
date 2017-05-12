/*
 * This is yaag middleware for the web apps using the middlewares that supports http handleFunc
 */
package middleware

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/masato25/yaag/yaag"
	"github.com/masato25/yaag/yaag/models"
)

/* 32 MB in memory max */
const MaxInMemoryMultipartSize = 32000000

var reqWriteExcludeHeaderDump = map[string]bool{
	"Host":              true, // not in Header map anyway
	"Content-Length":    true,
	"Transfer-Encoding": true,
	"Trailer":           true,
	"Accept-Encoding":   false,
	"Accept-Language":   false,
	"Cache-Control":     false,
	"Connection":        false,
	"Origin":            false,
	"User-Agent":        false,
}

type YaagHandler struct {
	nextHandler http.Handler
}

func Handle(nextHandler http.Handler) http.Handler {
	return &YaagHandler{nextHandler: nextHandler}
}

func (y *YaagHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !yaag.IsOn() {
		y.nextHandler.ServeHTTP(w, r)
		return
	}
	writer := httptest.NewRecorder()
	apiCall := models.ApiCall{}
	Before(&apiCall, r)
	y.nextHandler.ServeHTTP(writer, r)
	After(&apiCall, writer, w, r)
}

func HandleFunc(next func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if !yaag.IsOn() {
			next(w, r)
			return
		}
		apiCall := models.ApiCall{}
		writer := httptest.NewRecorder()
		Before(&apiCall, r)
		next(writer, r)
		After(&apiCall, writer, w, r)
	}
}

func Before(apiCall *models.ApiCall, req *http.Request) {
	apiCall.RequestHeader = ReadHeaders(req)
	apiCall.RequestUrlParams = ReadQueryParams(req)
	val, ok := apiCall.RequestHeader["Content-Type"]
	log.Println(val)
	if ok {
		ct := strings.TrimSpace(apiCall.RequestHeader["Content-Type"])
		switch ct {
		case "application/x-www-form-urlencoded":
			fallthrough
		case "application/json, application/x-www-form-urlencoded":
			log.Println("Reading form")
			apiCall.PostForm = ReadPostForm(req)
		case "application/json":
			log.Println("Reading body")
			apiCall.RequestBody = *ReadBody(req)
		default:
			if strings.Contains(ct, "multipart/form-data") {
				handleMultipart(apiCall, req)
			}
		}
	}
}

func ReadQueryParams(req *http.Request) map[string]string {
	params := map[string]string{}
	u, err := url.Parse(req.RequestURI)
	if err != nil {
		return params
	}
	for k, v := range u.Query() {
		if len(v) < 1 {
			continue
		}
		// TODO: v is a list, and we should be showing a list of values
		// rather than assuming a single value always, gotta change this
		params[k] = v[0]
	}
	return params
}

func printMap(m map[string]string) {
	for key, value := range m {
		log.Println(key, "=>", value)
	}
}

func handleMultipart(apiCall *models.ApiCall, req *http.Request) {
	apiCall.RequestHeader["Content-Type"] = "multipart/form-data"
	req.ParseMultipartForm(MaxInMemoryMultipartSize)
	apiCall.PostForm = ReadMultiPostForm(req.MultipartForm)
}

func ReadMultiPostForm(mpForm *multipart.Form) map[string]string {
	postForm := map[string]string{}
	for key, val := range mpForm.Value {
		postForm[key] = val[0]
	}
	return postForm
}

func ReadPostForm(req *http.Request) map[string]string {
	postForm := map[string]string{}
	log.Println("", *ReadBody(req))
	for _, param := range strings.Split(*ReadBody(req), "&") {
		value := strings.Split(param, "=")
		postForm[value[0]] = value[1]
	}
	return postForm
}

func ReadHeaders(req *http.Request) map[string]string {
	b := bytes.NewBuffer([]byte(""))
	err := req.Header.WriteSubset(b, reqWriteExcludeHeaderDump)
	if err != nil {
		return map[string]string{}
	}
	headers := map[string]string{}
	for _, header := range strings.Split(b.String(), "\n") {
		values := strings.Split(header, ":")
		if strings.EqualFold(values[0], "") {
			continue
		}
		headers[values[0]] = values[1]
	}
	return headers
}

func ReadHeadersFromResponse(writer *httptest.ResponseRecorder) map[string]string {
	headers := map[string]string{}
	for k, v := range writer.Header() {
		log.Println(k, v)
		headers[k] = strings.Join(v, " ")
	}
	return headers
}

func ReadBody(req *http.Request) *string {
	save := req.Body
	var err error
	if req.Body == nil {
		req.Body = nil
	} else {
		save, req.Body, err = drainBody(req.Body)
		if err != nil {
			return nil
		}
	}
	b := bytes.NewBuffer([]byte(""))
	chunked := len(req.TransferEncoding) > 0 && req.TransferEncoding[0] == "chunked"
	if req.Body == nil {
		return nil
	}
	var dest io.Writer = b
	if chunked {
		dest = httputil.NewChunkedWriter(dest)
	}
	_, err = io.Copy(dest, req.Body)
	if chunked {
		dest.(io.Closer).Close()
	}
	req.Body = save
	body := b.String()
	return &body
}

func After(apiCall *models.ApiCall, writer *httptest.ResponseRecorder, w http.ResponseWriter, r *http.Request) {
	if strings.Contains(r.RequestURI, ".ico") {
		fmt.Fprintf(w, writer.Body.String())
		return
	}
	if writer.Code != 404 {
		apiCall.MethodType = r.Method
		apiCall.CurrentPath = strings.Split(r.RequestURI, "?")[0]
		apiCall.ResponseBody = writer.Body.String()
		apiCall.ResponseCode = writer.Code
		apiCall.ResponseHeader = ReadHeadersFromResponse(writer)
		go yaag.GenerateHtml(apiCall)
	}
	for key, value := range apiCall.ResponseHeader {
		w.Header().Add(key, value)
	}
	w.WriteHeader(writer.Code)
	w.Write(writer.Body.Bytes())
}

// One of the copies, say from b to r2, could be avoided by using a more
// elaborate trick where the other copy is made during Request/Response.Write.
// This would complicate things too much, given that these functions are for
// debugging only.
func drainBody(b io.ReadCloser) (r1, r2 io.ReadCloser, err error) {
	var buf bytes.Buffer
	if _, err = buf.ReadFrom(b); err != nil {
		return nil, nil, err
	}
	if err = b.Close(); err != nil {
		return nil, nil, err
	}
	return ioutil.NopCloser(&buf), ioutil.NopCloser(bytes.NewReader(buf.Bytes())), nil
}
