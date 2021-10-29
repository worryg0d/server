package main

import (
	"fmt"
	"net/http"
	"os"
)

func PrintFile(w http.ResponseWriter, req *http.Request) {
	path := "." + req.URL.RequestURI()
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprint(w, err.Error())
		return
	}
	fmt.Fprint(w, string(data))
}

// HandleFunc For All Dirs Entries
func HandleFuncFADE(path string, f func(http.ResponseWriter, *http.Request)) {
	dirEnrtry, err := os.ReadDir(path)
	if err != nil {
		panic(err)
	}
	for _, val := range dirEnrtry {
		if val.IsDir() {
			HandleFuncFADE(path+"/"+val.Name(), f)
			continue
		}
		http.HandleFunc(path[1:]+"/"+val.Name(), f)
	}
}

func main() {
	HandleFuncFADE(".", PrintFile)
	http.ListenAndServe(":8080", nil)
}
