# webfingerprint
study project about wfp attacks on anonymous networks as tor jap etc

Requierments: 
1. Install libpcap-dev

2. Configure mysql server as a backend. Read database.txt 

3. Install selenium related things using tebeka's selenium package helper: 

```cd $GOPATH/src/github.com/tebeka/selenium/vendor/```

```go run init.go --alsologtostderr  --download_browsers --download_latest```

4. Install downloaded Chrome and Chrome-driver binaries in path. For example with creating symlinks from binaries to common $PATH folder under /usr/local/bin. 

```ln -s $GOPATH/src/github.com/tebeka/selenium/vendor/chrome-linux/chrome /usr/local/bin/chrome```

```ln -s $GOPATH/src/github.com/tebeka/selenium/vendor/chromedriver /usr/local/bin/chromedriver```

5. Build the binary with:

```go build -o bin/wfp cmd/*```

6. Set permissions for packet capture for built binary.

```sudo setcap cap_net_raw,cap_net_admin=eip bin/wfp```

7. ...

8. Profit! 
