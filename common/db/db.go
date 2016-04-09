package db

// graph.endpoint
type GraphEndpoint struct {
	Id       int64
	Endpoint string
}

// graph.tag_endpoint
type GraphTagEndpoint struct {
	Id         int64
	Tag        string
	EndpointId int64
}

// graph.endpoint_counter
type GraphEndpointCounter struct {
	Id         int64
	EndpointId int64
	Counter    string
}
