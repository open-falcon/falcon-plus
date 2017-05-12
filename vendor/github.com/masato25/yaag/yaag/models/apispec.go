package models

type ApiSpec struct {
	HttpVerb string
	Path     string
	Calls    []ApiCall
}
