package http

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	// "os"
)

const HOST = "eds.newtouch.cn"
const UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3359.181 Safari/537.36"
const AcceptLanguage = "zh-CN,zh;q=0.9,en;q=0.8,zh-TW;q=0.7"
const AcceptEncoding = "gzip, deflate"

var postProperties = make(map[string]string)
var getProperties = make(map[string]string)

func init() {
	log.Println("http init")

	// post
	postProperties["Host"] = HOST
	postProperties["Content-Length"] = "6955"
	postProperties["Cache-Control"] = "max-age=0"
	postProperties["Origin"] = "http://eds.newtouch.cn"
	postProperties["Upgrade-Insecure-Requests"] = "1"
	postProperties["Content-Type"] = "application/x-www-form-urlencoded"
	postProperties["User-Agent"] = "Mozilla/5.0"
	postProperties["Accept"] = "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3"
	postProperties["Accept-Encoding"] = AcceptEncoding
	postProperties["Accept-Language"] = AcceptLanguage
	postProperties["connection"] = "Keep-Alive"
	postProperties["accept"] = "*/*"
	postProperties["user-agent"] = "Mozilla/5.0"

	// get
	getProperties["Upgrade-Insecure-Requests"] = "1"
	getProperties["User-Agent"] = UserAgent
	getProperties["Accept"] = "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8"
	getProperties["Accept-Encoding"] = AcceptEncoding
	getProperties["Accept-Language"] = AcceptLanguage
	getProperties["Host"] = HOST
}

func DoGet(url string) (string, error) {
	return doRequest(url, http.MethodGet, nil)
}
func DoPost(url string, body io.Reader) (string, error) {
	return doRequest(url, http.MethodPost, body)
}

const cookie = "ASP.NET_SessionId=4khtnz55xiyhbmncrzmzyzzc; ActionSelect=010601; Hm_lvt_416c770ac83a9d996d7b3793f8c4994d=1569767826; Hm_lpvt_416c770ac83a9d996d7b3793f8c4994d=1569767826; PersonId=12234"

func doRequest(url, method string, body io.Reader) (string, error) {
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
	request.Header.Set("Cookie", cookie)

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
