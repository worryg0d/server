package main

import (
	"fmt"
	"net/http"
	"os"
)

func Explorer(path string) {
	dirEntry, err := os.ReadDir(path)
	if err != nil {
		panic(err)
	}
	for _, val := range dirEntry {
		if val.IsDir() {
			http.HandleFunc(path[1:]+"/"+val.Name()+"/", PrintDir)
			Explorer(path + "/" + val.Name())
		}
	}
}

func PrintDir(w http.ResponseWriter, req *http.Request) {
	path := "." + req.URL.RequestURI()
	dirEntry, err := os.ReadDir(path)
	if err != nil {
		panic(err)
	}
	fmt.Fprint(w, "<html style=\"font-size: 14pt\">")
	for _, val := range dirEntry {
		if val.IsDir() {
			fmt.Fprintf(w, "<a href=\"%s\">%s</a><br/>", path[1:]+val.Name()+"/", val.Name())
			continue
		}
		fmt.Fprintf(w, "<a href=\"%s\">%s</a><br/>", path[1:]+val.Name(), val.Name())
	}
	fmt.Fprintf(w, "</html>")
}

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
	http.HandleFunc("/", PrintDir)
	Explorer(".")
	http.ListenAndServe(":8080", nil)
}
