package restful

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SendOK(c *gin.Context, data interface{}) {
	c.JSON(200, gin.H{
		"code": "Success",
		"data": data,
		"msg":  "",
	})
}

func SendError(c *gin.Context, err error, data interface{}) {
	c.JSON(200, gin.H{
		"code": "Error",
		"data": data,
		"msg":  err.Error(),
	})
}

func DoPost(url string, body interface{}, params map[string]string, headers map[string]string) (*http.Response, error) {
	// add post body
	var bodyJson []byte
	var req *http.Request
	if body != nil {
		var err error
		bodyJson, err = json.Marshal(body)
		if err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(bodyJson))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-type", "application/json")
	// add params
	q := req.URL.Query()
	if params != nil {
		for key, val := range params {
			q.Add(key, val)
		}
		req.URL.RawQuery = q.Encode()
	}
	// add headers
	if headers != nil {
		for key, val := range headers {
			req.Header.Add(key, val)
		}
	}
	// http client
	client := &http.Client{}
	return client.Do(req)
}

// body, _ := ioutil.ReadAll(resp.Body)
func DoGet(url string, params map[string]string, headers map[string]string) (*http.Response, error) {
	// new request
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	// add params
	q := req.URL.Query()
	if params != nil {
		for key, val := range params {
			q.Add(key, val)
		}
		req.URL.RawQuery = q.Encode()
	}

	// add headers
	if headers != nil {
		for key, val := range headers {
			req.Header.Add(key, val)
		}
	}

	// http client
	client := &http.Client{}
	return client.Do(req)
}

// DoPut...
