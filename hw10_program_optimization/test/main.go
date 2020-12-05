package main

import (
	"bytes"
	"fmt"
	"github.com/dimazusov/hw-test/hw10_program_optimization"
	"log"
	"net/http"
	_ "net/http/pprof"
)

var data = `{"Id":1,"Name":"Howard Mendoza","Username":"0Oliver","Email":"aliquid_qui_ea@Browsedrive.gov","Phone":"6-866-899-36-79","Password":"InAQJvsq","Address":"Blackbird Place 25"}
{"Id":2,"Name":"Jesse Vasquez","Username":"qRichardson","Email":"mLynch@broWsecat.com","Phone":"9-373-949-64-00","Password":"SiZLeNSGn","Address":"Fulton Hill 80"}
{"Id":3,"Name":"Clarence Olson","Username":"RachelAdams","Email":"RoseSmith@Browsecat.com","Phone":"988-48-97","Password":"71kuz3gA5w","Address":"Monterey Park 39"}
{"Id":4,"Name":"Gregory Reid","Username":"tButler","Email":"5Moore@Teklist.net","Phone":"520-04-16","Password":"r639qLNu","Address":"Sunfield Park 20"}
{"Id":5,"Name":"Janice Rose","Username":"KeithHart","Email":"nulla@Linktype.com","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}`

func handler(w http.ResponseWriter, r *http.Request) {
	for i := 0; i < 10; i++ {
		hw10_program_optimization.GetDomainStat(bytes.NewBufferString(data), "com")
	}
	w.Write([]byte("ok"))
}

// go tool pprof -http=":8080" http://127.0.0.1:7070/debug/pprof/profile?second=5
func main() {
	fmt.Println("listen")
	http.HandleFunc("/", handler)

	err := http.ListenAndServe(":80", nil)
	log.Fatalln(err)
}
