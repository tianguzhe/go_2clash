## 交叉编译

```bash

 CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build main.go   
 CGO_ENABLED=1 GOOS=linux GOARCH=amd64 CC=/usr/local/gcc-4.8.1-for-linux64/bin/x86_64-pc-linux-gcc  go build main.go 

```

linux cgo的交叉编译

```bash

wget http://crossgcc.rts-software.org/download/gcc-4.8.1-for-linux32-linux64/gcc-4.8.1-for-linux64.dmg

export CGO_ENABLED=1
export GOOS=linux
export GOARCH=amd64
export CC=/usr/local/gcc-4.8.1-for-linux64/bin/x86_64-pc-linux-gcc
# 重新编译
go build
```

