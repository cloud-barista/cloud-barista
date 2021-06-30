package util

import (
	"encoding/base64"
	"fmt"

	// "reflect"
	// "io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httputil"

	// "net/url"
	"os"
	"strconv"
	"strings"

	// "time"
	"bytes"
	"encoding/json"
	"math"
	// "io/ioutil"
	// echosession "github.com/go-session/echo-session"
	// "github.com/labstack/echo"
	// "github.com/cloud-barista/cb-webtool/src/model"
)

type KeepZero float64

func (f KeepZero) MarshalJSON() ([]byte, error) {
	if float64(f) == float64(int(f)) {
		return []byte(strconv.FormatFloat(float64(f), 'f', 1, 32)), nil
	}
	return []byte(strconv.FormatFloat(float64(f), 'f', -1, 32)), nil
}

type myFloat64 float64

func (mf myFloat64) MarshalJSON() ([]byte, error) {
	const ε = 1e-12
	v := float64(mf)
	w, f := math.Modf(v)
	if f < ε {
		return []byte(fmt.Sprintf(`%v.0`, math.Trunc(w))), nil
	}
	return json.Marshal(v)
}

// ajax 호출할 때 header key 생성
func AuthenticationHandler() string {

	// conf 파일에 정의
	api_username := os.Getenv("API_USERNAME")
	api_password := os.Getenv("API_PASSWORD")
	// api_username := "default"
	// api_password := "default"

	//The header "KEY: VAL" is "Authorization: Basic {base64 encoded $USERNAME:$PASSWORD}".
	apiUserInfo := api_username + ":" + api_password
	encA := base64.StdEncoding.EncodeToString([]byte(apiUserInfo))
	//req.Header.Add("Authorization", "Basic"+encA)
	return "Basic " + encA

}

// originalUrl 은 API의 전체 경로
// parammapper 의 Key는 replace할 모든 text
// ex1) path인 경우 {abc}
// ex2) path인 경우 :abc
func MappingUrlParameter(originalUrl string, paramMapper map[string]string) string {
	returnUrl := originalUrl
	log.Println("originalUrl= ", originalUrl)
	if paramMapper != nil {
		for key, replaceValue := range paramMapper {
			returnUrl = strings.Replace(returnUrl, key, replaceValue, -1)
			// fmt.Println("Key:", key, "=>", "Element:", replaceValue+":"+returnUrl)
		}
	}
	log.Println("returnUrl= ", returnUrl)
	return returnUrl
}

// http 호출
func CommonHttp(url string, json []byte, httpMethod string) (*http.Response, error) {

	authInfo := AuthenticationHandler()

	log.Println("CommonHttp "+httpMethod+", ", url)
	// log.Println("authInfo ", authInfo)
	client := &http.Client{}
	req, err1 := http.NewRequest(httpMethod, url, bytes.NewBuffer(json))
	if err1 != nil {
		panic(err1)
	}

	// set the request header Content-Type for json
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	// req.Header.Set("Content-Type", "application/json")

	req.Header.Add("Authorization", authInfo)

	requestDump, err := httputil.DumpRequest(req, true)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(requestDump))
	resp, err := client.Do(req) // err 자체는 nil 이고 resp 내에 statusCode가 500임...

	return resp, err
}

// Json 형태의 bytes.Buffer 면 그대로 사용
func CommonHttpBytes(url string, jsonBytesBuffer *bytes.Buffer, httpMethod string) (*http.Response, error) {
	authInfo := AuthenticationHandler()

	// payload := strings.NewReader(`{
	// 	"CredentialName": "test-gcp-webtool21",
	// 	"ProviderName": "GCP",
	// 	"KeyValueInfoList": [
	// 		{
	// 			"Key": "project_id",
	// 			"Value": "megazone-for-yhnoh"
	// 		},
	// 		{
	// 			"Key": "client_email",
	// 			"Value": "yhnoh-704@megazone-for-yhnoh.iam.gserviceaccount.com"
	// 		},
	// 		{
	// 			"Key": "private_key",
	// 			"Value": "-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQDcbgpvhNXoYvq+\nVKtfcH/s0NL9Shmn8aSMd0eXGMMQ1b/VQ2HGUMxSgoOa0fHXOzIhpEKHtnIdv1uf\nad+AGmkQOUTRvmmcBUuh3JFuZpSoGH1InBZUowvzqaJvIrhPkAzuu+0el9kHdqRu\nwgCZaZXar0jXji3CvdTeAKtKnC1UgxTUrPAsfTiRh2ZkXDsIsrnA8kL8oIqEGm/g\ncx681KoRpbVb3LPt2BMXwZkEt9x1x+ExEBMjfrr6PEMkvzDlMrY8TcGSNXCHUAmN\no5o8Rsy88Rphd6viOxNuiQA7MacKBB3cR4fl12tfEzxbQUtGPFJB4doKE7AVb+EP\nCD3M6S1XAgMBAAECggEAUjBli7dH2uItBBKl42wbBr3GLdMXRdt/szA2bUw6T2ij\nom0BY+R0ir9HOs7VEZ9szcZlWBza5+SV0Ra00xsF2ZrA4kPRNO90h+GqCDQPca3P\n6ObqHJy+tBeoDTAw8NmROOKxQxrzPSkrnnCPsKQB0AxTaKwGu/n0COLO+37IGB4m\noOGtGIG+hPfLWGUegtZ9WM9m+widM+WXCeWmEORZ+k5hxE82XFH8ayWgj22FG+ys\n+kI8FjdZX4DNmSzrAl6v/KKlvPgqsEksSllRf/409XPGWSO/NK1th+3wsk69NS3Q\nbQzyZyWD99ubA50I5YAhpimmPFS3NGwOdiPKRa4moQKBgQD6CbsZAfK+Lf9nf14l\nrrKhPZTqNjOFpB+yPFgOhWgPVS+BVEXqcOOJ5H21hla6hYEJvhPAZcWGh9uq2xtX\nVieyuy3DhHXXHy3xaTnimTcxMR5KrD/4R/dB2S0qPmMdWDOZKYRN5PYA4PECtt+b\nURltqmC/spQxZ2BWVVlKfA+oGwKBgQDhr5PIulSZE7gZi1oefMV32vBlHJZrSUzi\n8TRVKoixYQ/0DcffSKdWpoaFcmoi+JjiBT5iKFYvzqwOEB6WuiVY9OHHqh6PeG+S\nQnE4UI1lyP9mmX8FWwXmuz8qghJjRKlHMW2engbOC6cTudYt0JWDWOeu4CmX+pZy\ncGjBLnubdQKBgFJXDgQoPhYe87LToN9r7mtm6jlO7BygdceuU5lEmYYjGWfPps6T\nqnrogfVbbggymtHohHyhhzDMYKydRx21w5D6TxHJ9zyGigysCGH07tYGROF2ZAKR\nQH2w1UzKCr3JJATWRTmZouGbMgMg0fZF+MfCieXXGzJBxtnndWYwAL/tAoGAKIyC\nUftgfcx0NGq8O2QRmrJEpPYY9JfL766Ex5SH0M7urdvYAH8uSbxLyShAd54Q4fMt\nPTegHKDWewRcappxYWVGN7iSGxb5fN7hNswKS7JsaQPFNbIgAk+8TqfmI93PSFJS\nLsCX2mdvknS+Tab/ZgUQQ3RVJNBKPa+CssrmPI0CgYEAv/Vi7HFjv0ozCaBR7cVm\n+n4J2L7nVw44tYuzhIYM1b82FU1OywGl7101m2TliJKV4nJEUaOvW4zUp1YEx+Uj\n2FkO/v2ZRRAqbebovw7+fium+Xjcipw7zjKTKJ9mAs87A3JABahvyUeGyZSw6xyX\no4fl0OvwIH9WjjayOEeriDk=\n-----END PRIVATE KEY-----\n"
	// 		}
	// 	]
	// }`)
	log.Println("CommonHttp "+httpMethod+", ", url)
	//[]byte(creds.PrivateKey)
	// log.Println("authInfo ", authInfo)
	client := &http.Client{}
	req, err1 := http.NewRequest(httpMethod, url, jsonBytesBuffer)
	// ppp := make([]byte, len(payload))
	// req, err1 := http.NewRequest(httpMethod, url, payload)
	if err1 != nil {
		panic(err1)
	}

	// url = "http://54.248.3.145:1323/tumblebug/ns/ns-01/resources/vNet"

	// set the request header Content-Type for json
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	// req.Header.Set("Content-Type", "application/json")

	req.Header.Add("Authorization", authInfo)

	// getBody := req.GetBody
	// copyBody, err := getBody()
	// if err == nil {
	// 	log.Println("------stert");
	// 	log.Println(copyBody);
	// 	log.Println("------end");
	// }
	requestDump, err := httputil.DumpRequest(req, true)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(requestDump))

	resp, err := client.Do(req) // err 자체는 nil 이고 resp 내에 statusCode가 500임...

	return resp, err
}

// func CommonHttpWithoutParam1(url string, httpMethod string) (io.ReadCloser, error) {
// 	authInfo := AuthenticationHandler()

// 	fmt.Println("CommonHttp ", url)
// 	client := &http.Client{}
// 	req, err := http.NewRequest(httpMethod, url, nil)
// 	if err != nil {
// 		panic(err)
// 	}

// 	// set the request header Content-Type for json
// 	req.Header.Set("Content-Type", "application/json; charset=utf-8")
// 	req.Header.Add("Authorization", authInfo)
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Println(resp.StatusCode)
// 	defer resp.Body.Close()

// 	return resp.Body, err
// }

// parameter 없이 호출하는 경우 사용.받은대로 return하면 호출하는 method에서 가공하여 사용
// func CommonHttpWithoutParam(url string, httpMethod string) (io.ReadCloser, error) {
// 	authInfo := AuthenticationHandler()

// 	fmt.Println("CommonHttp ", url)
// 	client := &http.Client{}
// 	req, err := http.NewRequest(httpMethod, url, nil)
// 	if err != nil {
// 		panic(err)
// 	}

// 	// set the request header Content-Type for json
// 	// req.Header.Set("Content-Type", "application/json; charset=utf-8")	// 사용에 주의할 것.
// 	req.Header.Add("Authorization", authInfo)
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		panic(err)
// 	}
// 	// respBody := resp.Body
// 	// robots, _ := ioutil.ReadAll(resp.Body)
// 	// defer resp.Body.Close()
// 	// log.Println(fmt.Print(string(robots)))
// 	// fmt.Println(resp.StatusCode)

// 	return resp.Body, err
// }

// parameter 없이 호출하는 경우 사용.받은대로 return하면 호출하는 method에서 가공하여 사용
func CommonHttpWithoutParam(url string, httpMethod string) (*http.Response, error) {
	authInfo := AuthenticationHandler()

	log.Println("CommonHttpWithoutParam "+httpMethod+", ", url)
	// log.Println("authInfo ", authInfo)
	client := &http.Client{}
	req, err := http.NewRequest(httpMethod, url, nil)
	if err != nil {
		fmt.Println("CommonHttpWithoutParam error")
		fmt.Println(err)
		panic(err)
	}

	// set the request header Content-Type for json
	// req.Header.Set("Content-Type", "application/json; charset=utf-8")	// 사용에 주의할 것.
	req.Header.Add("Authorization", authInfo)
	// resp, err := client.Do(req)
	return client.Do(req)
}

// Put/Post 등을 formData 형태로 호출할 때
// https://minwook-shin.github.io/go-decode-encode-url-values-form/ 참조할 것
//func CommonHttpFormData(targetUrl string, formParam url.Values, httpMethod string) (*http.Response, error) {
// func CommonHttpFormData(targetUrl string, formParam map[string]string, httpMethod string) (*http.Response, error) {
func CommonHttpFormData(targetUrl string, formParam map[string]interface{}, httpMethod string) (*http.Response, error) {
	//m := structs.Map(s)
	authInfo := AuthenticationHandler()

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)

	log.Println(formParam)

	// writer.WriteField 는 int 등으로는 전송이 안됨.... string으로 변환 후 전송
	for key, val := range formParam {
		_ = writer.WriteField(key, fmt.Sprintf("%v", val))
	}
	err := writer.Close()

	client := &http.Client{}
	req, _ := http.NewRequest(httpMethod, targetUrl, payload)

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Add("Authorization", authInfo)

	resp, err := client.Do(req)
	return resp, err
}

// return message 확인용
func DisplayResponse(resp *http.Response) {
	fmt.Println("*****DisplayResponse begin****")
	if resp == nil {
		log.Println(" response is nil ")
	} else {
		// resultMessage, err1 := ioutil.ReadAll(resp.Message)
		// if err1 != nil {
		// 	str := string(resultMessage)
		// 	println("nil ", str)
		// 	println("err1 ", err1)
		// }
		// fmt.Println(string(resultMessage))
		// log.Println(" 11111111111111111111111111111 ")

		fmt.Println(resp.StatusCode)
		log.Println(" 22222222222222222222222222 ")

		fmt.Println(string(resp.Status))
		log.Println(" 3333333333333333333 ")
		// data, err := ioutil.ReadAll(resp.Body)
		// if err != nil {
		//     panic(err)
		// }
		// fmt.Printf("%s\n", string(data))

		// resultStatus := resp.StatusCode
		// fmt.Println("resultStatus ", resultStatus)
		// // fmt.Println("body ",  resp.Body)
		resultBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			str := string(resultBody)
			println("nil ", str)
			println("err ", err)
		}
		fmt.Println(string(resultBody))

		var target interface{}
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal(body, &target)
		fmt.Println(fmt.Println(target))
		// // json.NewDecoder(respBody).Decode(&stringMap)
		// pbytes, _ := json.Marshal(resultBody)
		// fmt.Println(string(pbytes))

		fmt.Println("*****DisplayResponse end****")
	}
}

// Response 객체의 내용
// type Response struct {
//     Status     string // e.g. "200 OK"
//     StatusCode int    // e.g. 200
//     Proto      string // e.g. "HTTP/1.0"
//     ProtoMajor int    // e.g. 1
//     ProtoMinor int    // e.g. 0

//     // response headers
//     Header http.Header
//     // response body
//     Body io.ReadCloser
//     // request that was sent to obtain the response
//     Request *http.Request
// }
