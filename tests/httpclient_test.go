package tests

import (
	"github.com/liangboceo/yuanboot/pkg/httpclient"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetHttp(t *testing.T) {
	httpServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		body, _ := ioutil.ReadAll(r.Body)

		if r.Header.Get("Content-Type") == "application/json" {
			if string(body) != `{"aid":1002,"auth":"ss"}` {
				t.Error("RunPost json type func error")
			}
		}

		if r.Header.Get("Content-Type") == "application/x-www-form-urlencoded" {
			if string(body) != `word=你好` {
				t.Error("RunPost WithFormRequest type func error 'word=你好' != ", string(body))
			}
		}

		if r.Header.Get("Content-Type") == "application/text" {
			w.WriteHeader(200)
			_, _ = w.Write([]byte("hello"))
		}
	}))
	defer httpServer.Close()

	c := httpclient.NewClient()
	request := httpclient.WithJsonRequest(`{"aid":1002,"auth":"ss"}`).POST(httpServer.URL)
	_, _ = c.Do(request)

	request1 := httpclient.WithFormRequest(map[string]interface{}{
		"word": "你好",
	}).POST(httpServer.URL)
	_, _ = c.Do(request1)

	request2 := httpclient.WithRequest().Header("Content-Type", "application/text").GET(httpServer.URL)
	resp, _ := c.Do(request2)

	assert.Equal(t, resp.GetRequestTime().Seconds() < 5, true)
	assert.Equal(t, string(resp.Body), "hello")
}

//func TestHttpClientFactoryBaseUrl(t *testing.T) {
//	httpServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		w.WriteHeader(200)
//		_, _ = w.Write([]byte("ok"))
//	}))
//	fmt.Print(httpServer.URL)
//	//httpServer.URL = "http://127.0.0.1:8080"
//	defer httpServer.Close()
//	factory := httpclient.NewFactory()
//	client, err := factory.Create(httpServer.URL)
//	if err != nil {
//		panic(err)
//	}
//	req := httpclient.WithRequest()
//	req.GET("")
//	req.SetTimeout(10)
//	res, err := client.Do(req)
//	fmt.Print(string(res.Body))
//	assert.Equal(t, string(res.Body), "ok")
//}

func TestBaseUrlSplicingUrl(t *testing.T) {
	factory := httpclient.NewFactory()
	client1, err := factory.Create("http")
	if err == nil {
		panic("fail")
	}
	_, err2 := factory.Create("http://www.baidu.com")
	if err2 != nil {
		panic("fail")
	}
	_, err3 := factory.Create("http/www.baidu.com")
	if err3 == nil {
		panic("fail")
	}
	_, err4 := factory.Create("http://c.cn")
	if err4 != nil {
		panic("fail")
	}

	res := client1.SplicingUrl("http://www.baidu.com", "/v1/user")
	assert.Equal(t, "http://www.baidu.com/v1/user", res)

	res = client1.SplicingUrl("http://www.baidu.com/", "/v1/user")
	assert.Equal(t, "http://www.baidu.com/v1/user", res)

	res = client1.SplicingUrl("http://www.baidu.com/", "/v1/user/")
	assert.Equal(t, "http://www.baidu.com/v1/user", res)

}
