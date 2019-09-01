# punc

Profile functions(*P*rofile f*UNC*tions) for golang.

```go
package main

import (
        "log"
        "time"

        "github.com/tkuchiki/punc"
        "github.com/tkuchiki/punc/httpserver"
)

func trace() {
        defer punc.Done(punc.Do())
}

func main() {
        go func() {
        	
                log.Println(httpserver.ListenAndServe())
        }()

        time.Sleep(time.Second * 2)

        trace()
        trace()
        trace()
 
        time.Sleep(time.Second * 3600)
}
```

```console
# terminal 1
# default: PUNC_HOST="localhost:58080"
$ go run main.go

# terminal 2
$ curl -s "localhost:58080/stats" > stats.csv

$ cat stats.csv
count,func,call_stack,max,min,sum,avg,p50,p99
3,trace,main.main>main.trace,0.000002,0.000000,0.000003,0.000001,0.000002,0.000001
```
