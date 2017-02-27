# gorotatewriter

Sorry for poor english.
This file can rotate log file by time.Duration.
Checking for rotating performing at any write operation, so it can have high cost for performance

#example
```go

package main

import (
	"log"
	"io"
	"os"
	"time"
	"fmt"
)

func main(){
    rotWriter, err := New("./log.log", 24 * time.Hour) // or maybe 10 * time.Minute

    if err!=nil{
      fmt.Println("unable init writer")
      return
    }

    multi := io.MultiWriter(rotWriter, os.Stdout)
    logger := log.New(multi, "", log.Ldate|log.Ltime)
    
    logger.Println("test")
}
```
