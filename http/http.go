package http

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	// "os"
)

const UrlHome = "http://eds.newtouch.cn"
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
	postProperties["Origin"] = UrlHome
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

// DoRequest http请求，返回 string
func DoRequest(url, method, cookie string, body io.Reader) string {
	// log.Println(url)
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		panic(err)
	}
	request.Header.Set("Referer", url)

	var headerProps map[string]string
	if method == http.MethodPost {
		headerProps = postProperties
	} else {
		headerProps = getProperties
	}
	for key, value := range headerProps {
		request.Header.Set(key, value)
	}

	request.Header.Set("Cookie", cookie)

	resp, err := http.DefaultClient.Do(request)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	// 查看一下
	// io.Copy(os.Stdout, resp.Body)

	// log.Println(resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		//io.Copy(os.Stdout, resp.Body)
		log.Panicln(resp.StatusCode)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	return string(respBody)
}
