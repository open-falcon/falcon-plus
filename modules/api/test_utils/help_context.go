package thelp

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
)

//https://github.com/gin-gonic/gin/issues/580#issuecomment-242168245
func NewTestContext(method, path string, body *[]byte) (w *httptest.ResponseRecorder, r *http.Request) {
	w = httptest.NewRecorder()
	if body == nil {
		r, _ = http.NewRequest(method, path, bytes.NewBuffer(nil))
	} else {
		r, _ = http.NewRequest(method, path, bytes.NewBuffer(*body))
	}
	//set json post as default
	r.Header.Set("Content-Type", "application/json")
	r.PostForm = url.Values{}
	return
}

func NewTestContextWithDefaultSession(method, path string, body *[]byte) (w *httptest.ResponseRecorder, r *http.Request) {
	w, r = NewTestContext(method, path, body)
	r = setSession(r)
	return w, r
}

func setSession(r *http.Request) *http.Request {
	r = SetSessionWith(r, "testuser92", "48380ba36ad211e79fb3001500c6ca5a")
	return r
}

func SetDefaultAdminSession(r *http.Request) *http.Request {
	r = SetSessionWith(r, "root", "4833cd596ad211e79fb3001500c6ca5a")
	return r
}
func SetSessionWith(r *http.Request, name string, sig string) *http.Request {
	r.Header.Set("Apitoken", fmt.Sprintf("{\"name\":\"%v\",\"sig\":\"%v\"}", name, sig))
	return r
}

func CleanSession(r *http.Request) *http.Request {
	r.Header.Set("Apitoken", "")
	return r
}
