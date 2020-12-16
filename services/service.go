package services

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/yann0917/dedao-dl/request"
	"github.com/yann0917/dedao-dl/utils"
)

var (
	dedaoCommURL = &url.URL{
		Scheme: "https",
		Host:   "dedao.cn",
	}
	baseURL = "https://www.dedao.cn"
)

// Response dedao success response
type Response struct {
	H respH `json:"h"`
	C respC `json:"c"`
}

type respH struct {
	C int    `json:"c"`
	E string `json:"e"`
	S int    `json:"s"`
	T int    `json:"t"`
}

// respC response content
type respC []byte

func (r *respC) UnmarshalJSON(data []byte) error {
	*r = data

	return nil
}

func (r respC) String() string {
	return string(r)
}

//Service dedao service
type Service struct {
	client *request.HTTPClient
}

//NewService new service
func NewService(gat, isid, sid, acwTc, iget, token, guardDeviceID string) *Service {
	client := request.NewClient(baseURL)
	client.ResetCookieJar()
	cookies := []*http.Cookie{}
	cookies = append(cookies, &http.Cookie{
		Name:   "GAT",
		Value:  gat,
		Domain: "." + dedaoCommURL.Host,
	})
	cookies = append(cookies, &http.Cookie{
		Name:   "ISID",
		Value:  isid,
		Domain: "." + dedaoCommURL.Host,
	})
	cookies = append(cookies, &http.Cookie{
		Name:   "_guard_device_id",
		Value:  guardDeviceID,
		Domain: "www." + dedaoCommURL.Host,
	})
	cookies = append(cookies, &http.Cookie{
		Name:   "_sid",
		Value:  sid,
		Domain: "www." + dedaoCommURL.Host,
	})
	cookies = append(cookies, &http.Cookie{
		Name:   "acw_tc",
		Value:  acwTc,
		Domain: "www." + dedaoCommURL.Host,
	})
	cookies = append(cookies, &http.Cookie{
		Name:   "iget",
		Value:  iget,
		Domain: "www." + dedaoCommURL.Host,
	})
	cookies = append(cookies, &http.Cookie{
		Name:   "token",
		Value:  token,
		Domain: "www." + dedaoCommURL.Host,
	})
	client.Client.Jar.SetCookies(dedaoCommURL, cookies)

	return &Service{client: client}
}

//Cookies get cookies string
func (s *Service) Cookies() map[string]string {
	cookies := s.client.Client.Jar.Cookies(dedaoCommURL)

	cstr := map[string]string{}

	for _, cookie := range cookies {
		cstr[cookie.Name] = cookie.Value
	}

	return cstr
}

func (r *Response) isSuccess() bool {
	return r.H.C == 0
}

func deferResponseClose(s *http.Response) {
	if s != nil {
		defer s.Body.Close()
	}
}

func handleHTTPResponse(resp *http.Response, err error) (io.ReadCloser, error) {
	if err != nil {
		deferResponseClose(resp)
		return nil, err
	}
	if resp.StatusCode == 452 {
		return nil, errors.New("452")
	}
	return resp.Body, nil
}

func handleJSONParse(reader io.Reader, v interface{}) error {
	result := new(Response)

	err := utils.UnmarshalReader(reader, &result)
	if err != nil {
		fmt.Printf("err1: %s \n", err.Error())
		return err
	}
	// fmt.Printf("result.C:=%#v", result.C)
	if !result.isSuccess() {
		//未登录或者登录凭证无效
	}
	err = utils.UnmarshalJSON(result.C, v)
	if err != nil {
		fmt.Printf("err2: %s", err.Error())
		return err
	}

	return nil
}