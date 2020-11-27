package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
)

const (
	defaultPort       = ":8080"
	defaultCarpetURL  = "http://europe-west3-playground-296517.cloudfunctions.net/carpet?degree=5"
	defaultHTMLFormat = `%s`
)

type config struct {
	serverPort string
	carpetURL  string
	htmlFormat string
}

func main() {
	logrus.Println("Get config")
	cfg := getConfig()

	logrus.Println("Build handler")
	router := mux.NewRouter()
	router.HandleFunc("/", buildHandler(cfg.carpetURL, cfg.htmlFormat))
	router.Use(mux.CORSMethodMiddleware(router))

	logrus.Println(fmt.Sprintf("Start listening on port: %s", cfg.serverPort))
	http.ListenAndServe(cfg.serverPort, router)
}

func getConfig() config {
	cfg := config{}
	cfg.carpetURL = os.Getenv("CARPET_URL")
	if cfg.carpetURL == "" {
		cfg.carpetURL = defaultCarpetURL
	}

	cfg.serverPort = os.Getenv("SERVER_ADDRESS")
	if cfg.serverPort == "" {
		cfg.serverPort = defaultPort
	}

	htmlFormat, err := ioutil.ReadFile("index.html.format")
	if err != nil {
		cfg.htmlFormat = defaultHTMLFormat
	} else {
		cfg.htmlFormat = string(htmlFormat)
	}

	return cfg
}

func buildHandler(carpetURL, htmlFormat string) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		resp, err := http.Get(carpetURL)
		if err != nil {
			logrus.Fatalln(fmt.Sprintf("get error: %s", err.Error()))
			http.Error(writer, err.Error(), 500)
			return
		}
		defer resp.Body.Close()

		carpet, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logrus.Fatalln(fmt.Sprintf("read file error: %s", err.Error()))
			http.Error(writer, err.Error(), 500)
			return
		}

		logrus.Println("Carpet generated")
		fmt.Fprintf(writer, htmlFormat, string(carpet))
	}
}
