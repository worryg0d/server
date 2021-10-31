package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func explorer(path string) error {
	dirEntry, err := os.ReadDir(path)
	if err != nil {
		log.Printf("cannot read dir on path: %s. Error: %s", path, err.Error())
		return err
	}
	for _, val := range dirEntry {
		if val.IsDir() {
			http.HandleFunc(path[1:]+"/"+val.Name()+"/", printDir)
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
			fmt.Sprintf("cannot read dir on (Path: %s). rr: %s", path, err.Error()),
			http.StatusInternalServerError,
		)
		log.Printf("cannot read dir (Path: %s), err: %s", path, err.Error())
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
		if val.IsDir() {
			registerFilesHandlers(path+"/"+val.Name(), f)
			continue
		}
		http.HandleFunc(path[1:]+"/"+val.Name(), f)
	}

	return nil
}

func main() {

	err := registerFilesHandlers(".", printFile)
	if err != nil {
		log.Fatalln(err)
	}

	http.HandleFunc("/", printDir)

	err = explorer(".")
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Server is starting on port: 8080")
	err = http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Fatalf("cannot start server. Error: %s", err.Error())
	}
}
