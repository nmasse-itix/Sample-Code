package main
 
import (
	"os"
	"fmt"
	"flag"
        "net/http"
        "net/http/httputil"
        "net/url"
)

// The local port on which we should listen to
var local_port int 

// The target URL on which we should redirect requests 
var target string 

func init() {
	// --local-port=9009
	flag.IntVar(&local_port, "local-port", 9090, "the TCP port to listen to")
	
	// --target=http://www.perdu.com
	flag.StringVar(&target, "target", "http://www.perdu.com", "the target URL to redirect the request to")
}

// The MyResponseWriter is a wrapper around the standard http.ResponseWriter
// We need it to retrieve the http status code of an http response
type MyResponseWriter struct {
	Underlying	http.ResponseWriter
	Status		int
}

func (mrw *MyResponseWriter) Header() http.Header {
	return mrw.Underlying.Header()
}

func (mrw *MyResponseWriter) Write(b []byte) (int, error) {
	return mrw.Underlying.Write(b)
}

func (mrw *MyResponseWriter) WriteHeader(s int) {
	mrw.Status = s
	mrw.Underlying.WriteHeader(s)
} 

func main() {
	// Parse the command line arguments
	flag.Parse()

	// Parse the target URL and perform some sanity checks
	url, err := url.Parse(target)
	if err != nil {
		panic(err)
	}
	
	// Initialize a Reverse Proxy object with a custom director
        proxy := httputil.NewSingleHostReverseProxy(url)
        underlying_director := proxy.Director
	proxy.Director = func(req *http.Request) {
		// Let the underlying director do the mandatory job
		underlying_director(req)

		// Custom Handling
		// ---------------
		//
		// Filter out the "Host" header sent by the client 
		// otherwise the target server won't be able to find the 
		// matching virtual host. The correct host header will be
		// added automatically by the net/http package.
        	req.Host = ""
        }

	http.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		// Log the incoming request (including headers)
		fmt.Printf("%v %v HTTP/1.1\n", req.Method, req.URL)
		req.Header.Write(os.Stdout)
		fmt.Println()

		// Wrap the standard response writer with our own 
		// implementation because we need the status code of the 
		// response and that field is not exported by default
		mrw := &MyResponseWriter{ Underlying: rw }

		// Let the reverse proxy handle the request
		proxy.ServeHTTP(mrw, req)
		
		// Log the response
		fmt.Printf("%v %v\n", mrw.Status, http.StatusText(mrw.Status))
		mrw.Header().Write(os.Stdout)
		fmt.Println()
	})
	
	fmt.Printf("Listening on port %v for incoming requests...\n", local_port)
        http.ListenAndServe(fmt.Sprintf(":%v", local_port), nil)
}

