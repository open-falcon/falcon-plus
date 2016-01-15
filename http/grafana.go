package http

import (
	"bytes"
	"encoding/json"
	"github.com/open-falcon/query/g"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Idc struct {
	Id         int
	Pop_id     int
	Name       string
	Count      int
	Area       string
	Province   string
	City       string
	Updated_at string
}

type Province struct {
	Id         int
	Province   string
	Count      int
	Updated_at string
}

type City struct {
	Id         int
	City       string
	Province   string
	Count      int
	Updated_at string
}

func getHosts(w http.ResponseWriter, req *http.Request, hostKeyword string) {
	if len(hostKeyword) == 1 {
		hostKeyword = ".+"
	}
	rand.Seed(time.Now().UTC().UnixNano())
	random64 := rand.Float64()
	_r := strconv.FormatFloat(random64, 'f', -1, 32)
	maxQuery := strconv.Itoa(g.Config().Api.Max)
	url := "/api/endpoints" + "?q=" + hostKeyword + "&tags&limit=" + maxQuery + "&_r=" + _r + "&regex_query=1"
	if strings.Index(g.Config().Api.Query, req.Host) >= 0 {
		url = "http://localhost:9966" + url
	} else {
		url = g.Config().Api.Query + url
	}

	reqGet, err := http.NewRequest("GET", url, nil)
	if err != nil {
		StdRender(w, "", err)
	}

	client := &http.Client{}
	resp, err := client.Do(reqGet)
	if err != nil {
		StdRender(w, "", err)
	}

	defer resp.Body.Close()

	result := []interface{}{}
	if resp.Status == "200 OK" {
		body, _ := ioutil.ReadAll(resp.Body)
		var nodes = make(map[string]interface{})
		if err := json.Unmarshal(body, &nodes); err != nil {
			StdRender(w, "", err)
		}
		for _, host := range nodes["data"].([]interface{}) {
			item := map[string]interface{}{
				"text":       host,
				"expandable": true,
			}
			result = append(result, item)
		}
		RenderJson(w, result)
	} else {
		RenderJson(w, result)
	}
}

func getNextCounterSegment(metric string, counter string) string {
	if len(metric) > 0 {
		metric += "."
	}
	counter = strings.Replace(counter, metric, "", 1)
	segment := strings.Split(counter, ".")[0]
	return segment
}

func checkSegmentExpandable(segment string, counter string) bool {
	segments := strings.Split(counter, ".")
	expandable := !(segment == segments[len(segments)-1])
	return expandable
}

func getMetrics(w http.ResponseWriter, req *http.Request, query string) {
	result := []interface{}{}

	query = strings.Replace(query, "#.*", "", -1)
	arrQuery := strings.Split(query, "#")
	host, arrMetric := arrQuery[0], arrQuery[1:]
	maxQuery := strconv.Itoa(g.Config().Api.Max)
	metric := strings.Join(arrMetric, ".")
	reg, _ := regexp.Compile("(^{|}$)")
	host = reg.ReplaceAllString(host, "")
	host = strings.Replace(host, ",", "\",\"", -1)

	endpoints := "[\"" + host + "\"]"

	rand.Seed(time.Now().UTC().UnixNano())
	random64 := rand.Float64()
	_r := strconv.FormatFloat(random64, 'f', -1, 32)

	form := url.Values{}
	form.Set("endpoints", endpoints)
	form.Add("q", metric)
	form.Add("limit", maxQuery)
	form.Add("_r", _r)

	target := "/api/counters"
	if strings.Index(g.Config().Api.Query, req.Host) >= 0 {
		target = "http://localhost:9966" + target
	} else {
		target = g.Config().Api.Query + target
	}

	reqPost, err := http.NewRequest("POST", target, strings.NewReader(form.Encode()))
	if err != nil {
		log.Println("Error =", err.Error())
	}
	reqPost.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(reqPost)
	if err != nil {
		log.Println("Error =", err.Error())
	}
	defer resp.Body.Close()

	if resp.Status == "200 OK" {
		body, _ := ioutil.ReadAll(resp.Body)
		var nodes = make(map[string]interface{})
		if err := json.Unmarshal(body, &nodes); err != nil {
			log.Println(err.Error())
		}
		var segmentPool = make(map[string]int)
		for _, data := range nodes["data"].([]interface{}) {
			counter := data.([]interface{})[0].(string)
			segment := getNextCounterSegment(metric, counter)
			expandable := checkSegmentExpandable(segment, counter)
			if _, ok := segmentPool[segment]; !ok {
				item := map[string]interface{}{
					"text":       segment,
					"expandable": expandable,
				}
				result = append(result, item)
				segmentPool[segment] = 1
			}
		}
		RenderJson(w, result)
	} else {
		RenderJson(w, result)
	}
}

func setQueryEditor(w http.ResponseWriter, req *http.Request) {
	query := req.URL.Query().Get("query")
	query = strings.Replace(query, ".%", "", -1)
	query = strings.Replace(query, ".undefined", "", -1)
	query = strings.Replace(query, ".select metric", "", -1)
	if !strings.Contains(query, "#") {
		getHosts(w, req, query)
	} else {
		getMetrics(w, req, query)
	}
}

func getMetricValues(req *http.Request, host string, targets []string, result []interface{}) []interface{} {
	endpoint_counters := []interface{}{}
	metric := strings.Join(targets, ".")
	if strings.Contains(host, "{") { // Templating metrics request
		host = strings.Replace(host, "{", "", -1)
		host = strings.Replace(host, "}", "", -1)
		hosts := strings.Split(host, ",")
		for _, host := range hosts {
			item := map[string]string{
				"endpoint": host,
				"counter":  metric,
			}
			endpoint_counters = append(endpoint_counters, item)
		}
	} else {
		item := map[string]string{
			"endpoint": host,
			"counter":  metric,
		}
		endpoint_counters = append(endpoint_counters, item)
	}

	if len(endpoint_counters) > 0 {
		from, err := strconv.ParseInt(req.PostForm["from"][0], 10, 64)
		until, err := strconv.ParseInt(req.PostForm["until"][0], 10, 64)
		url := "/graph/history"
		if strings.Index(g.Config().Api.Query, req.Host) >= 0 {
			url = "http://localhost:9966" + url
		} else {
			url = g.Config().Api.Query + url
		}

		args := map[string]interface{}{
			"start":             from,
			"end":               until,
			"cf":                "AVERAGE",
			"endpoint_counters": endpoint_counters,
		}
		bs, err := json.Marshal(args)
		if err != nil {
			log.Println("Error =", err.Error())
		}

		reqPost, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(bs)))
		if err != nil {
			log.Println("Error =", err.Error())
		}
		reqPost.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(reqPost)
		if err != nil {
			log.Println("Error =", err.Error())
		}
		defer resp.Body.Close()

		if resp.Status == "200 OK" {
			body, _ := ioutil.ReadAll(resp.Body)
			nodes := []interface{}{}
			if err := json.Unmarshal(body, &nodes); err != nil {
				log.Println(err.Error())
			}

			for _, node := range nodes {
				if _, ok := node.(map[string]interface{})["Values"]; ok {
					result = append(result, node)
				}
			}
		}
	}
	return result
}

func getValues(w http.ResponseWriter, req *http.Request) {
	result := []interface{}{}
	req.ParseForm()
	for _, target := range req.PostForm["target"] {
		if !strings.Contains(target, ".select metric") {
			targets := strings.Split(target, "#")
			host, targets := targets[0], targets[1:]
			result = getMetricValues(req, host, targets, result)
		}
	}
	RenderJson(w, result)
}

func GrafanaApiParser(w http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		setQueryEditor(w, req)
	} else if req.Method == "POST" {
		getValues(w, req)
	}
}

func configGrafanaRoutes() {
	http.HandleFunc("/api/grafana/", GrafanaApiParser)
}
