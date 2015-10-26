package restcl

import "net/http"

//rest object, used to built and query endpoints
type Rest struct {
	Prefix      string
	interceptor []Intercept
	endpoints   map[string]RestEndPoint
	tmp         RestEndPoint
}

type Intercept func (req *http.Request)


//creates a new rest object
func NewRest() *Rest {
	b := &Rest{}
	b.endpoints = make(map[string]RestEndPoint)
	b.interceptor = nil
	b.tmp = RestEndPoint{}
	return b
}

//adds a prefix to the url's
func (r *Rest) SetPrefix(prefix string) *Rest {
	r.Prefix = prefix
	return r
}

//adds a intercept interface to the interceptor slice
func (r *Rest) Use(inter Intercept) *Rest {
	if r.tmp.Method == "" {
		r.interceptor = append(r.interceptor, inter)
	} else {
		r.tmp.interceptor = append(r.tmp.interceptor, inter)
	}
	return r
}

//creates a new endpoint
func (r *Rest) Create(url string) *Rest {
	r.tmp = RestEndPoint{}
	r.tmp.Url = r.Prefix + url
	if r.interceptor != nil {
		if len(r.tmp.interceptor) == 0 {
			r.tmp.interceptor = r.interceptor
		} else {
			for _, inte := range r.interceptor {
				r.tmp.interceptor = append(r.tmp.interceptor, inte)
			}
		}
	}
	return r
}

//sets the method of the endpoint
func (r *Rest) SetMethod(method string) *Rest {
	r.tmp.Method = method
	return r
}

//adds the endpoint to the endpoints slice
func (r *Rest) Build(name string) *Rest {
	r.endpoints[name] = r.tmp
	r.tmp = RestEndPoint{}
	return r
}

//returns the endpoint
func (r *Rest) Get(name string) *RestEndPoint {
	rep := r.endpoints[name]
	return &rep
}
