package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const CLIENT_NUM = 3

var (
	API_KEY = ""
	ips     [CLIENT_NUM]string
)

func main() {
	API_KEY = os.Getenv("API_KEY")

	http.HandleFunc("/", handleRequest)

	fmt.Printf("Starting server...\n")
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
		if r.Header.Get("Authorization") != "API_KEY: "+API_KEY {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("403 - Forbidden"))
			return
		}
		result, err := getIP()
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(result))
	case "POST":
		if r.Header.Get("Authorization") != "API_KEY: "+API_KEY {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("403 - Forbidden"))
			return
		}
		if err := r.ParseForm(); err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		err := storeIP(r.FormValue("id"), r.RemoteAddr)
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
	ip_array := strings.Split(ip, ":")
	id_num, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	if id_num >= CLIENT_NUM {
		return errors.New("id outside of client range")
	}
	ips[id_num] = ip_array[0]
	return nil
}

func getIP() (string, error) {
	jsonData, err := json.Marshal(ips)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}
