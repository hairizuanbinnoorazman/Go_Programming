package main

import (
	"log"
	"net/http"
	"text/template"
)

type basicWebsite struct{}

func (b *basicWebsite) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./templates/header.tmpl",
		"./templates/index.tmpl",
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}
	err = ts.ExecuteTemplate(w, "index", nil)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}

type GoogleAnalyticsParameters struct {
	// General
	ProtocolVersion string `json:"protocol_version"`
	TrackingID      string `json:"tracking_id"`
	// User
	ClientID string `json:"client_id"`
	// Content Information
	DocumentLocationURL string `json:"document_location_url"`
	// System Info
	ScreenResolution         string `json:"screen_resolution"`
	ViewportSize             string `json:"viewport_size"`
	UserLanguage             string `json:"user_language"`
	UserAgentArchitecture    string `json:"user_agent_architecture"`
	UserAgentFullVersionList string `json:"user_agent_full_version_list"`
	UserAgentMobile          bool   `json:"user_agent_mobile"`
	UserAgentModel           string `json:"user_agent_model"`
	UserAgentPlatform        string `json:"user_agent_platform"`
	UserAgentPlatformVersion string `json:"user_agent_platform_version"`
	// Hit
	HitType           string `json:"hit_type"`
	NonInteractionHit bool   `json:"non_interaction_hit"`
}

type analytics struct{}

// example request:
// http://localhost:8080/analytics/collect?v=1&_v=j99&a=1160670874&t=pageview&_s=1&dl=http%3A%2F%2Flocalhost%2Findex&ul=en-us&de=UTF-8&sd=24-bit&sr=1920x1080&vp=1920x487&je=0&_u=QACAAUABAAAAAAAAII~&jid=&gjid=&cid=329073308.1678561437&tid=UA-93097338-1&_gid=1043476193.1678561437&_fplc=0&gtm=457e3360&z=1118969917
// Reference:
// https://www.thyngster.com/ga4-measurement-protocol-cheatsheet/
func (a *analytics) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("start processing analytics request")
	defer log.Println("end processing analytics request")

	ga_params := GoogleAnalyticsParameters{}
	// General
	ga_params.ProtocolVersion = r.URL.Query().Get("v")
	ga_params.TrackingID = r.URL.Query().Get("tid")
	// User
	ga_params.ClientID = r.URL.Query().Get("cid")
	// Content Information
	ga_params.DocumentLocationURL = r.URL.Query().Get("dl")
	// System Info
	ga_params.ScreenResolution = r.URL.Query().Get("sr")
	ga_params.ViewportSize = r.URL.Query().Get("vp")
	ga_params.UserLanguage = r.URL.Query().Get("ul")
	ga_params.UserAgentArchitecture = r.URL.Query().Get("uaa")
	ga_params.UserAgentFullVersionList = r.URL.Query().Get("uafvl")
	if r.URL.Query().Get("uamb") == "1" {
		ga_params.UserAgentMobile = true
	}
	ga_params.UserAgentModel = r.URL.Query().Get("uam")
	ga_params.UserAgentPlatform = r.URL.Query().Get("uap")
	ga_params.UserAgentPlatformVersion = r.URL.Query().Get("uapv")
	// Hit
	ga_params.HitType = r.URL.Query().Get("t")
	if r.URL.Query().Get("ni") == "1" {
		ga_params.NonInteractionHit = true
	}

	log.Printf("%+v\n", ga_params)

}

func main() {
	http.Handle("/index", &basicWebsite{})
	http.Handle("/analytics/collect", &analytics{})
	http.Handle("/analytics/g/collect", &analytics{})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
