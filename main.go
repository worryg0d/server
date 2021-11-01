package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

var handlersList = make(map[string]struct{})

func explorer(path string) error {
	dirEntry, err := os.ReadDir(path)
	if err != nil {
		log.Printf("cannot read dir on path: %s. Error: %s", path, err.Error())
		return err
	}
	for _, val := range dirEntry {
		if val.IsDir() {
			p := path+"/"+val.Name()+"/"
			if _, ok := handlersList[p]; ok {
				continue
			}
			handlersList[p] = struct{}{}
			http.HandleFunc(p[1:]+"/", printDir)
			explorer(path + "/" + val.Name())
		}
	}
	return nil
}

func printDir(w http.ResponseWriter, req *http.Request) {
	path := "." + req.URL.RequestURI()
	dirEntries, err := os.ReadDir(path)

	if err != nil {
		http.Error(w,
			fmt.Sprintf("cannot print dir on (Path: %s), err: %s", path, err.Error()),
			http.StatusInternalServerError,
		)
		log.Printf("cannot print dir (Path: %s), err: %s", path, err.Error())
		return
	}

	fmt.Fprint(w, "<html style=\"font-size: 14pt\">")

	for _, val := range dirEntries {
		if val.IsDir() {
			fmt.Fprintf(w, "<a href=\"%s\">%s (folder)</a><br/>", path[1:]+val.Name()+"/", val.Name())
			continue
		}
		fmt.Fprintf(w, "<a href=\"%s\">%s</a><br/>", path[1:]+val.Name(), val.Name())
	}

	fmt.Fprintf(w, "</html>")
}

func printFile(w http.ResponseWriter, req *http.Request) {
	path := "." + req.URL.RequestURI()

	data, err := os.ReadFile(path)

	if err != nil {
		http.Error(w,
			fmt.Sprintf("cannot read file (Path: %s). Error: %s", path, err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	fmt.Fprint(w, string(data))
}

func registerFilesHandlers(path string, f func(http.ResponseWriter, *http.Request)) error {
	dirEntries, err := os.ReadDir(path)
	if err != nil {
		log.Printf("cannot read dir (Path: %s). Error: %s", path, err.Error())
		return err
	}
	for _, val := range dirEntries {
		p := path+"/"+val.Name()
		if val.IsDir() {
			registerFilesHandlers(p, f)
			continue
		}
		if _, ok := handlersList[p]; ok {
			continue
		}
		handlersList[p] = struct{}{}

		http.HandleFunc(p[1:], f)
	}

	return nil
}

func registerHandlersInRealTime(path string)  {
	for {
		err := registerFilesHandlers(path, printFile)
		if err != nil {
			log.Fatalln(err)
		}

		err = explorer(path)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func main() {

	http.HandleFunc("/", printDir)
	go registerHandlersInRealTime(".")
	log.Printf("Server is starting on port: 8080")
	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Fatalf("cannot start server. Error: %s", err.Error())
	}
}
