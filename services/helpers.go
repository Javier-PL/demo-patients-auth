package services

import (
	"bytes"
	"encoding/json"
	"io/ioutil"

	"net/http"
	"time"
)

var CustomClient = &http.Client{Timeout: 120 * time.Second}

//var Log log.Logger

type ReqData struct {
	Method string
	URL    string
	Token  string
	Body   interface{}
	Target interface{}
}

type ResData struct {
	Code    int
	ResBody []byte
	Error   error
}

func (d *ReqData) DoReq() (r *ResData) {
	var bodyAsBytes []byte

	if d.Body != nil {
		var marshalErr error
		bodyAsBytes, marshalErr = json.Marshal(d.Body)
		if marshalErr != nil {
			Log.Error("ERROR: Cannot convert body to JSON")
			return &ResData{Code: 500, ResBody: nil, Error: marshalErr}
		}
	} else {
		bodyAsBytes = nil
	}

	req, err := http.NewRequest(d.Method, d.URL, bytes.NewBuffer(bodyAsBytes))

	//fmt.Println(req)
	if d.Token != "" {
		req.Header.Add("Authorization", d.Token)
		//fmt.Println("TESEEST", "Zoho-oauthtoken "+token)
		//fmt.Println("TESEEST", "Zoho-oauthtoken "+os.Getenv("ZOHO_TOKEN"))
		//req.Header.Add("Authorization", token)
	}

	res, err := CustomClient.Do(req)
	if err != nil {
		Log.Error("ERROR: Cannot " + req.Method + " JSON to: " + d.URL)
		return &ResData{Code: 500, ResBody: nil, Error: err} //TODO: here a goroutine failed. Error was runtime error: invalid memory address or nil pointer dereference. I suppose that there was no statuscode.Maybe I shoudl return statuscode as pointer or return 500
	}
	body, _ := ioutil.ReadAll(res.Body)

	//fmt.Println(body)
	defer res.Body.Close()

	Log.Debug(res.StatusCode, req.Method, "> from > ", d.URL)

	if d.Target != nil {
		errDecode := json.Unmarshal([]byte(string(body)), &d.Target)
		return &ResData{Code: res.StatusCode, ResBody: body, Error: errDecode}
	}

	return &ResData{Code: res.StatusCode, ResBody: body, Error: nil}
}
