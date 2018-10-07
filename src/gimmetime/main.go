package main

import (
	"fmt"
	"github.com/gocraft/web"
	"log"
	"net/http"
	"time"
)

//dummy contex for web request state
type Context struct {
}

//returns current time as number of sec from 1.1.1970
func TimeNow(rw web.ResponseWriter, req *web.Request) {
	log.Println("TimeNow called...")
	rw.Write([]byte(fmt.Sprintf("%v", time.Now().Unix())))
}

func main() {

	router := web.New(Context{}).Middleware(web.ShowErrorsMiddleware)
	router.Get("/now", TimeNow)

	log.Println("Service gimmetime started and listening on port 80")
	http.ListenAndServe("0.0.0.0:80", router)

}
