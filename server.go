package main

import (
	"html/template"
	"net/http"
	"log"
	"fmt"
	"io"
	"os"
	"strings"
	"encoding/json"
)

func handleIndex(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/index.html")
	if (err != nil) {
		return
	}

	t.Execute(w, nil)
}


type WeatherResponse struct {
	Features []struct{
		Properties struct {
			Name string `json:"name"`
		} `json:"properties`
	} `json:"features"`
}

func handleZone(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	zoneList, err := http.Get(fmt.Sprintf("https://api.weather.gov/zones/forecast?area=%s", strings.ToUpper(state)))
	if err != nil {
        fmt.Print(err.Error())
        os.Exit(1)
    }

	var response WeatherResponse
	defer zoneList.Body.Close()

	var resList []string

	if zoneList.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(zoneList.Body)
		if err != nil {
			log.Fatal(err)
		}
		err = json.Unmarshal(bodyBytes, &response)
		if err != nil {
			fmt.Println("error unmarshalling: ", err)
		}


		for _, feature := range response.Features {
			resList = append(resList, feature.Properties.Name)
		}
	}


	t, err := template.ParseFiles("templates/zones.html")
	if err != nil {
		log.Fatal(err)
	}

	data := struct {
		State string
		Results []string
	  }{
		State: state,
		Results: resList,
	  }

	t.Execute(w, data)
}

func main() {
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/zones", handleZone)

	err := http.ListenAndServe(":3000", nil)
	if (err != nil) {
		return
	}
}