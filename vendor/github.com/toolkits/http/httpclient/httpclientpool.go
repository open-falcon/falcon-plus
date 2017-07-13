package httpclient

import (
	"github.com/toolkits/container/nmap"
	"net/http"
	"time"
)

//
var (
	httpClientPool = NewHttpClientPool()
)

func GetHttpClient(name string, connTimeout time.Duration, reqTimeout time.Duration) *http.Client {
	return httpClientPool.AddAndGetHttpClient(name, connTimeout, reqTimeout)
}

func RemoveHttpClient(name string) {
	httpClientPool.RemoveHttpClient(name)
}

// HttpClientPool Struct
type HttpClientPool struct {
	httpClientMap *nmap.SafeMap
}

func NewHttpClientPool() *HttpClientPool {
	return &HttpClientPool{httpClientMap: nmap.NewSafeMap()}
}

func (hcp *HttpClientPool) AddHttpClient(name string, connTimeout time.Duration, reqTimeout time.Duration) *http.Client {
	hci, found := hcp.httpClientMap.Get(name)
	if found {
		return hci.(*http.Client)
	}

	nhc := hcp.newHttpClient(connTimeout, reqTimeout)
	hcp.httpClientMap.Put(name, nhc)
	return nhc
}

func (hcp *HttpClientPool) GetHttpClient(name string) (*http.Client, bool) {
	hci, found := hcp.httpClientMap.Get(name)
	if found {
		return hci.(*http.Client), true
	}
	return &http.Client{}, false
}

func (hcp *HttpClientPool) AddAndGetHttpClient(name string, connTimeout time.Duration, reqTimeout time.Duration) *http.Client {
	hci, found := hcp.httpClientMap.Get(name)
	if found {
		return hci.(*http.Client)
	}

	nhc := hcp.newHttpClient(connTimeout, reqTimeout)
	hcp.httpClientMap.Put(name, nhc)
	return nhc
}

func (hcp *HttpClientPool) Size() int {
	return hcp.httpClientMap.Size()
}

func (hcp *HttpClientPool) RemoveHttpClient(name string) {
	if client, found := hcp.httpClientMap.GetAndRemove(name); found {
		if client.(*http.Client).Transport != nil {
			client.(*http.Client).Transport.(*Transport).Close()
		}
	}
}

func (hcp *HttpClientPool) RemoveAllHttpClients() {
	for _, key := range hcp.httpClientMap.Keys() {
		clienti, found := hcp.httpClientMap.GetAndRemove(key)
		if found {
			client := clienti.(*http.Client)
			if client.Transport != nil {
				client.Transport.(*Transport).Close()
			}
		}
	}
}

// internal
func (hcp *HttpClientPool) newHttpClient(connTimeout time.Duration, reqTimeout time.Duration) *http.Client {
	transport := &Transport{
		ConnectTimeout: connTimeout,
		RequestTimeout: reqTimeout,
	}
	return &http.Client{Transport: transport}
}
