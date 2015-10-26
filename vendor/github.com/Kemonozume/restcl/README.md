#restcl

restcl is a simple to use rest client

## Getting Started

Install restcl
~~~  go
go get github.com/Kemonozume/restcl
~~~ 

Start Using it
~~~ go
package main

import (
	"fmt"
	"github.com/Kemonozume/restcl"
)

type Post struct {
	UserId int
	Id     int
	Title  string
	Body   string
}

func main() {
	//define the rest structure
	rest := restcl.NewRest()
	rest.SetPrefix("http://jsonplaceholder.typicode.com")
	rest.Create("/posts").SetMethod("GET").Build("posts")
	rest.Create("/posts").SetMethod("POST").Build("createpost")
	rest.Create("/posts/{id}").SetMethod("GET").Build("post")
	rest.Create("/posts/{id}").SetMethod("DELETE").Build("deletepost")

	//get the post
	var post Post
	err := rest.Get("post").SetParams("id", 1).Exec(&post)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", post)

	//create a post
	var post1 Post
	err = rest.Get("createpost").
		SetBody("userId", 1, "title", "test", "body", "test").
		Exec(&post1)

	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", post1)
	/* will return
	{
	  id: 101,
	}
	*/
}
~~~
modify the request using the Intercept interface
~~~ go
package main

import (
	"fmt"
	"net/http"
	"github.com/Kemonozume/restcl"
)

type Auth struct{}
func (a Auth) Modify(req *http.Request) {
	//do something with the request
}

type Log struct{}
func (l Log) Modify(req *http.Request) {
	//log the request
}

func main() {
	//define the rest structure
	rest := restcl.NewRest()
	//logs every request but only uses the auth intercept on post/delete
	rest.SetPrefix("http://jsonplaceholder.typicode.com").Use(Log{})
	rest.Create("/posts").SetMethod("GET").Build("posts")
	rest.Create("/posts").SetMethod("POST").Use(Auth{}).Build("createpost")
	rest.Create("/posts/{id}").SetMethod("GET").Build("post")
	rest.Create("/posts/{id}").SetMethod("DELETE").Use(Auth{}).Build("deletepost")
}
~~~


## Contributing
Feel free to put up a Pull Request.

## About

restcl is multi-threading safe to use and easy to customize 

project is inspired by [gorilla/mux](github.com/gorilla/mux)  
