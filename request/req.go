package request

import (
	"io/ioutil"
	"log"
	"net/http"
)

/**
构建请求下载
*/
type ReqH map[string]string
type Requests struct {
	Url     string
	Method  string
	Headers ReqH
}

func (req Requests) News() (*http.Response,[]byte) {

	NewReq, _ := http.NewRequest(req.Method, req.Url, nil)
	for key, value := range req.Headers {
		NewReq.Header.Add(key, value)
	}
	client := &http.Client{}
	resp, err := client.Do(NewReq)
	if err !=nil {
		log.Printf("地址：%s 请求失败",req.Url)

	}
	defer resp.Body.Close()
	data, _ := ioutil.ReadAll(resp.Body)
	return resp, data
}
