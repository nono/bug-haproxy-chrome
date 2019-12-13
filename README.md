Bug haproxy - Chrome
====================

Description
-----------

I work for [Cozy Coud](https://cozy.io/) and we have found an issue in
production, with a weird `net:ERR_HTTP2_PROTOCOL_ERROR` error. After a bit of
debugging, we have identified that it happens when
Chrome/Chromium/Electron/NodeJS uploads a file in HTTP/2 to haproxy, that
proxies this request to our backend in Go (HTTP/1.1 between haproxy and the
backend), and that the backend sends a 409 Conflict response before the client
has sent the whole file. I'm not sure if it is a bug in haproxy, a bug in
Chrome&co, or just me that does something I shouldn't. In this repository, I've
tried to make it simple to reproduce the bug.


Steps to reproduce
------------------

1. Start the go server in a terminal with `go run main.go`
2. Configure and start haproxy:

```
mkdir -p ssl/bug.localhost
openssl genrsa -out ssl/bug.localhost/bug.localhost.key 1024
openssl req -new -key ssl/bug.localhost/bug.localhost.key -out ssl/bug.localhost/bug.localhost.csr
openssl x509 -req -days 365 -in ssl/bug.localhost/bug.localhost.csr -signkey ssl/bug.localhost/bug.localhost.key -out ssl/bug.localhost/bug.localhost.crt
cat ssl/bug.localhost/bug.localhost.crt ssl/bug.localhost/bug.localhost.key > ssl/bug.localhost/bug.localhost.pem
echo "127.0.0.1 bug.localhost" | sudo tee -a /etc/hosts
haproxy -f haproxy.conf
```

3. Open chrome on the index page with `chromium-browser https://localhost:1443/`
   (and ignore the security warning)
4. Open the devtools in chrome, and select the network tab
5. On the web form, select a big file

What is expected: Chrome receive a 409 response.

What happens: Chrome show in the console the error
`POST https://localhost:1443/upload net::ERR_HTTP2_PROTOCOL_ERROR`.

Versions
--------

I've tried with Chrome 79 and Chromium 78. For haproxy, I've used the latest
stable (2.1.1), and some older versions (1.9.13 in particular), compiled with
`make -j 8 TARGET=linux-glibc USE_OPENSSL=1 USE_ZLIB=1 USE_PCRE=1`.
