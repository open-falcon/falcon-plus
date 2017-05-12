package models

type ApiCall struct {
	Id int

	CurrentPath string
	MethodType  string

	PostForm map[string]string

	RequestHeader        map[string]string
	CommonRequestHeaders map[string]string
	ResponseHeader       map[string]string
	RequestUrlParams     map[string]string

	RequestBody  string
	ResponseBody string
	ResponseCode int
}
