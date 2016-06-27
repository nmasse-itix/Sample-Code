## Reverse Proxy written in GO
This reverse proxy listens on a local port and forward all requests to a named host. It honors the proxy environment variables. 

### Initial Need

Once upon a time, I had to circumvent a bug in a product that could not handle correctly an HTTPS connection to a proxy. 

Since it was an HTTPS connection, I could not setup a transparent proxy. 

### What it does

It opens a local port and listen to HTTP requests, forwards the requests to a named host and send back the response. 

### How to use it

```bash
go run src/itix.fr/forward/main.go -local-port 8080 -target https://www.opentrust.com 
curl -D - http://localhost:8080/robots.txt
``` 

If you want to go through a proxy, do not forget to set the ```http_proxy``` and ```https_proxy``` variables !

```bash
export http_proxy=http://my.proxy:8888/
export https_proxy=http://my.proxy:8888/
``` 

