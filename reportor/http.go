package reportor

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
	// "os"
)

const UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3359.181 Safari/537.36"
const AcceptLanguage = "zh-CN,zh;q=0.9,en;q=0.8,zh-TW;q=0.7"
const AcceptEncoding = "gzip, deflate"

var postProperties = map[string]string{
	"Content-Length":            "6955",
	"Cache-Control":             "max-age=0",
	"Upgrade-Insecure-Requests": "1",
	"Content-Type":              "application/x-www-form-urlencoded",
	"User-Agent":                "Mozilla/5.0",
	"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3",
	"Accept-Encoding":           AcceptEncoding,
	"Accept-Language":           AcceptLanguage,
	"connection":                "Keep-Alive",
	"accept":                    "*/*",
	"user-agent":                "Mozilla/5.0",
}

var getProperties = map[string]string{
	"Upgrade-Insecure-Requests": "1",
	"User-Agent":                UserAgent,
	"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8",
	"Accept-Encoding":           AcceptEncoding,
	"Accept-Language":           AcceptLanguage,
}

const cookie = "ASP.NET_SessionId=4khtnz55xiyhbmncrzmzyzzc; ActionSelect=010601; Hm_lvt_416c770ac83a9d996d7b3793f8c4994d=1569767826; Hm_lpvt_416c770ac83a9d996d7b3793f8c4994d=1569767826; PersonId=12234"

func newClient() *http.Client {
	return &http.Client{
		Timeout: 30 * time.Second,
	}
}

func (r WorkReportor) doGet(url string) (string, error) {
	return r.doRequest(url, http.MethodGet, nil)
}
func (r WorkReportor) doPost(url string, body io.Reader) (string, error) {
	return r.doRequest(url, http.MethodPost, body)
}

func (r WorkReportor) doRequest(url, method string, body io.Reader) (string, error) {
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return "", err
	}

	var headerProps map[string]string
	if method == http.MethodPost {
		headerProps = postProperties
	} else {
		headerProps = getProperties
	}
	for key, value := range headerProps {
		request.Header.Set(key, value)
	}

	request.Header.Set("Referer", url)
	// cookie 不能省，否则做其他请求报错，如获取隐藏参数
	request.Header.Set("Cookie", cookie)
	//request.Header.Set("Host", e.host)
	//request.Header.Set("Origin", e.urls["home"])

	resp, err := newClient().Do(request)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 查看一下
	// io.Copy(os.Stdout, resp.Body)

	// log.Println(resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		msg, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("cannot read body: %w", err)
		}
		return "", fmt.Errorf("%w: %s, %s",
			err, http.StatusText(resp.StatusCode), msg)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(respBody), nil
}
