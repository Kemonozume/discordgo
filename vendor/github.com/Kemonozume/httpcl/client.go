package httpcl

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type Client struct {
	Error      error
	StatusCode int
	client     *http.Client
	redirect   bool
	request    *http.Request
}

type ClientBuilder struct {
	Method   string
	Url      string
	Redirect bool
	Body     []interface{}
}

func (c ClientBuilder) Build() *Client {
	cl := getRequestWithBody(c.Method, c.Url, c.Body)
	cl.redirect = c.Redirect
	return cl
}

//creates a http client using GET
func Get(url string) *Client {
	c := &Client{}
	c.request, c.Error = http.NewRequest("GET", url, nil)
	c.redirect = true
	return c
}

//creates a http client using HEAD
func Head(url string) *Client {
	c := &Client{}
	c.request, c.Error = http.NewRequest("HEAD", url, nil)
	c.redirect = true
	return c
}

//creates a http client using DElETE
func Delete(url string) *Client {
	c := &Client{}
	c.request, c.Error = http.NewRequest("DELETE", url, nil)
	c.redirect = true
	return c
}

//creates a http client using PATCH with the given params
func Patch(purl string, params ...interface{}) *Client {
	return getRequestWithBody("PATCH", purl, params)
}

//creates a http client using POST with the given params
func Post(purl string, params ...interface{}) *Client {
	return getRequestWithBody("POST", purl, params)
}

//creates a http client using PUT with the given params
func Put(purl string, params ...interface{}) *Client {
	return getRequestWithBody("PUT", purl, params)
}

//creates a request with a body
func getRequestWithBody(method, purl string, params []interface{}) *Client {
	if len(params) == 1 {
		switch params[0].(type) {
		case map[string]interface{}:
			return postMap(method, purl, params[0].(map[string]interface{}))
		case io.Reader:
			c := &Client{}
			c.request, c.Error = http.NewRequest(method, purl, params[0].(io.Reader))
			return c
		case url.Values:
			c := &Client{}
			c.request, c.Error = http.NewRequest(method, purl, strings.NewReader(params[0].(url.Values).Encode()))
			return c
		default:
			c := &Client{}
			c.Error = errors.New(fmt.Sprintf("parameters not correct %T", params[0]))
			return c
		}
	} else if len(params) > 1 {
		arr := params
		if len(arr)%2 == 0 {
			c := &Client{}
			values := url.Values{}
			for i := 0; i < len(arr)-1; i += 2 {
				c.Error = addToPost(arr[i].(string), arr[i+1], &values)
			}
			if c.Error != nil {
				return c
			}
			c.request, c.Error = http.NewRequest(method, purl, strings.NewReader(values.Encode()))
			return c
		} else {
			c := &Client{}
			c.Error = errors.New(fmt.Sprintf("parameters not correct, expected %v parameter got %v", len(arr)+1, len(arr)))
			return c
		}
	}
	c := &Client{}
	c.request, c.Error = http.NewRequest(method, purl, nil)
	return c
}

func addToPost(key string, value interface{}, values *url.Values) (err error) {
	switch value.(type) {
	case bool:
		values.Add(key, strconv.FormatBool(value.(bool)))
	case float64:
		values.Add(key, strconv.FormatFloat(value.(float64), 'f', 6, 64))
	case int:
		values.Add(key, strconv.Itoa(value.(int)))
	case int64:
		values.Add(key, strconv.FormatInt(value.(int64), 10))
	case rune:
		values.Add(key, string(value.(rune)))
	case string:
		values.Add(key, value.(string))
	case uint64:
		values.Add(key, strconv.FormatUint(value.(uint64), 10))
	default:
		err = errors.New(fmt.Sprintf("unsupported type for post %T", value))
	}
	return
}

func postMap(method, purl string, param map[string]interface{}) *Client {
	c := &Client{}
	values := url.Values{}
	for key, value := range param {
		c.Error = addToPost(key, value, &values)
	}
	if c.Error != nil {
		return c
	}
	c.request, c.Error = http.NewRequest(method, purl, strings.NewReader(values.Encode()))
	return c
}

//stops the httpclient from following redirects
func redirect(req *http.Request, via []*http.Request) error {
	return errors.New("no redirect")
}

//returns error "no request" if the client has no request
func (c *Client) hasRequest() error {
	if c.request == nil {
		return errors.New("no request")
	}
	return nil
}

//executes the function if the client has a request or not
func (c *Client) runWithHasRequest(s func()) *Client {
	if err := c.hasRequest(); err != nil {
		c.Error = err
		return c
	} else {
		s()
		return c
	}
}

//returns the underlying http.Request
func (c *Client) GetRequest() *http.Request {
	return c.request
}

//sets the underlying http.Request
func (c *Client) SetRequest(req *http.Request) *Client {
	c.request = req
	return c
}

//returns the underlying http.Client
func (c *Client) GetClient() *http.Client {
	return c.client
}

//sets the underlying http.Client
func (c *Client) SetClient(cl *http.Client) *Client {
	c.client = cl
	return c
}

//adds a header to the request
func (c *Client) AddHeader(key string, value string) *Client {
	return c.runWithHasRequest(func() {
		c.request.Header.Add(key, value)
	})
}

//adds header to the request using a map[string]string
func (c *Client) AddHeaderMap(header map[string]string) *Client {
	return c.runWithHasRequest(func() {
		for key, value := range header {
			c.request.Header.Add(key, value)
		}
	})
}

//sets the user agent for the request
func (c *Client) SetUserAgent(value string) *Client {
	return c.runWithHasRequest(func() {
		c.request.Header.Add("User-Agent", value)
	})
}

func (c *Client) SetBasicAuth(user, password string) *Client {
	return c.runWithHasRequest(func() {
		c.request.SetBasicAuth(user, password)
	})
}

//adds a cookie to the request
func (c *Client) AddCookie(cookie *http.Cookie) *Client {
	return c.runWithHasRequest(func() {
		c.request.AddCookie(cookie)
	})
}

//adds a cookies slice to the request
func (c *Client) AddCookies(cookies []*http.Cookie) *Client {
	return c.runWithHasRequest(func() {
		for _, cookie := range cookies {
			c.request.AddCookie(cookie)
		}
	})
}

//specify if you want to follow redirects or not
func (c *Client) FollowRedirect(red bool) *Client {
	return c.runWithHasRequest(func() {
		c.redirect = red
	})
}

//starts the request
func (c *Client) Do() (*http.Response, error) {
	if err := c.hasRequest(); err == nil {
		if c.Error != nil {
			return nil, c.Error
		} else {
			if c.client == nil {
				if c.redirect {
					c.client = &http.Client{}
				} else {
					c.client = &http.Client{
						CheckRedirect: redirect,
					}
				}
			}
			resp, err := c.client.Do(c.request)
			if err != nil {
				if !c.redirect {
					if !strings.Contains(err.Error(), "no redirect") {
						c.Error = err
					}
				} else {
					c.Error = err
				}
			}
			if resp != nil {
				c.StatusCode = resp.StatusCode
			} else {
				c.StatusCode = -1
			}
			return resp, c.Error
		}
	} else {
		c.Error = err
		return nil, err
	}
}

//starts the request and transforms the response with the given function
func (c *Client) DoTransform(trans func(resp *http.Response, c interface{}) error, b interface{}) (resp *http.Response, err error) {
	resp, err = c.Do()
	if err != nil {
		return nil, err
	}

	return resp, trans(resp, b)
}

//simple json transform
func TransformToJson(resp *http.Response, c interface{}) (err error) {
	defer resp.Body.Close()
	by, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(by, c)
	return
}

//simple string transform
func TransformToString(resp *http.Response, c interface{}) (err error) {
	defer resp.Body.Close()
	str, ok := c.(*string)
	if !ok {
		return errors.New(fmt.Sprintf("expected *string, got %T", c))
	}
	by, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	*str = string(by)
	return
}
