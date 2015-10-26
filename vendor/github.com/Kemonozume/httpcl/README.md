#httpcl

httpcl is a simple, chainable wrapper around the stdlib http client

## Getting Started

Install httpcl
~~~  go
go get github.com/Kemonozume/httpcl
~~~ 

Start Using it
~~~ go
package main

import (
    "github.com/Kemonozume/httpcl"
    "fmt"
)

func main() {
  resp, err := httpcl.Get("http://www.google.de").Do()
  if err != nil {
    panic(err)
  }
  fmt.Printf("status code: %d\n", resp.StatusCode) 
}
~~~

##Examples

POST 

supported parameters are io.Reader, url.Values, map[string]interface{} or key,value pairs 

key,value pairs and map[string]inteface{} have limited type support
(bool, float64, int, int64, rune, string, uint64)
~~~ go
package main

import (
	"fmt"

	"github.com/Kemonozume/httpcl"
)

func main() {
	//using simple key,value pairs
	//gets encoded by url.Values
	resp, err := httpcl.Post("http://httpbin.org/post", "test", "===value").Do()
	if err != nil {
		panic(err)
	}
	fmt.Println(resp.Status)

	//using a map
	vals := map[string]interface{}{
		"test":   "value",
		"number": 3000,
	}
	resp, err = httpcl.Post("http://httpbin.org/post", vals).Do()
	if err != nil {
		panic(err)
	}
	fmt.Println(resp.Status)
}
~~~

Transform response directly using helper functions
~~~ go
package main

import (
	"fmt"
	"github.com/Kemonozume/httpcl"
)

type UserAgent struct {
	Name string `json:"user-agent"`
}

func main() {
	var agent UserAgent
	var str string

	cl := httpcl.Get("http://httpbin.org/user-agent").
		SetUserAgent("httpcl")

	//each DoTransform is a seperate http request
	cl.DoTransform(httpcl.TransformToJson, &agent)
	cl.DoTransform(httpcl.TransformToString, &str)

	fmt.Printf("%s\n", agent.Name)
	fmt.Println(str)
}
~~~
use your own functions
~~~ go
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"github.com/Kemonozume/httpcl"
)

type UserAgent struct {
	Name string `json:"user-agent"`
}

func trans_to_useragent(resp *http.Response, c interface{}) (err error) {
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return errors.New("status not 200")
	}
	user := c.(*UserAgent)
	by, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(by, user)
	return
}

func main() {
	var agent UserAgent

	err := httpcl.Get("http://httpbin.org/user-agent").
		SetUserAgent("httpcl").
		DoTransform(trans_to_useragent, &agent)

	if err != nil {
		panic(err)
	}

	fmt.Println(agent.Name)

}
~~~


## Contributing
Feel free to put up a Pull Request.

## About

My first go library inspired by [jcabi http](http://http.jcabi.com/index.html)

[coverage](https://dl.dropboxusercontent.com/u/17033881/coverage.html)
