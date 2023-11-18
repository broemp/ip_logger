package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

const client_num = 3

var ips [client_num]string

func main() {
	http.HandleFunc("/", handleRequest)

	fmt.Printf("Starting server for testing HTTP POST...\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	switch r.Method {
	case "GET":
		result, err := getIP()
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(result))
	case "POST":
		if err := r.ParseForm(); err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		err := storeIP(r.FormValue("id"), r.FormValue("ip"))
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		w.Write([]byte("Success"))
	default:
		io.WriteString(w, "Error: only POST and GET allowed")
	}
}

func storeIP(id string, ip string) error {
	id_num, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	if id_num >= client_num {
		return errors.New("id outside of client range")
	}
	ips[id_num] = ip
	return nil
}

func getIP() (string, error) {
	jsonData, err := json.Marshal(ips)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}
