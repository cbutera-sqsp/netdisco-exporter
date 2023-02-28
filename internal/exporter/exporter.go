package exporter

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	log "github.com/sirupsen/logrus"
)

var (
	netdiscoHost     = os.Getenv("NETDISCO_HOST")
	netdiscoUsername = os.Getenv("NETDISCO_USERNAME")
	netdiscoPassword = os.Getenv("NETDISCO_PASSWORD")

	// Prom metrics
	NETDISCO_API_STATUS = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "NETDISCO_API_STATUS",
		Help: "Netdisco API status",
	})

	NETDISCO_LAST_DISCOVER = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "NETDISCO_LAST_DISCOVER",
		Help: "Netdisco last discover",
	},
		[]string{"hostname"},
	)

	NETDISCO_LAST_ARPNIP = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "NETDISCO_LAST_ARPNIP",
		Help: "Netdisco last arpnip",
	},
		[]string{"hostname"},
	)

	NETDISCO_LAST_MACSUCK = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "NETDISCO_LAST_MACSUCK",
		Help: "Netdisco last macsuck",
	},
		[]string{"hostname"},
	)
)

// Check API Status
func ApiStatus() {
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	path := "/swagger.json"
	url := netdiscoHost + path
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Error(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		NETDISCO_API_STATUS.Set(0)
		log.Error(err)
	} else if resp.StatusCode == 200 {
		NETDISCO_API_STATUS.Set(1)
		defer resp.Body.Close()
	}
}

// Login to netdisco API to retrieve API key for subsequent requests
func Login() string {
	path := "/login"
	url := netdiscoHost + path
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	req.SetBasicAuth(netdiscoUsername, netdiscoPassword)
	req.Header.Add("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	Apikey := Key{}
	err = json.Unmarshal(body, &Apikey)
	if err != nil {
		log.Fatal(err)
	}
	return Apikey.Data
}

// Logout to netdisco API - Destroy user API Key and session cookie
func Logout(ApiKey string) {
	path := "/logout"
	url := netdiscoHost + path
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", ApiKey)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Error(resp.StatusCode)
		log.Error("Logout successful")
	}
}

// Struct containing API key
type Key struct {
	Data string `json:"api_key"`
}

// Data from API path /api/v1/search/device
type Device []struct {
	Location          string      `json:"location"`
	LastMacsuckStamp  interface{} `json:"last_macsuck_stamp"`
	DNS               string      `json:"dns"`
	Model             string      `json:"model"`
	SinceLastArpnip   float64     `json:"since_last_arpnip"`
	UptimeAge         string      `json:"uptime_age"`
	FirstSeenStamp    string      `json:"first_seen_stamp"`
	LastArpnipStamp   string      `json:"last_arpnip_stamp"`
	SinceLastMacsuck  float64     `json:"since_last_macsuck"`
	Serial            string      `json:"serial"`
	LastDiscoverStamp string      `json:"last_discover_stamp"`
	SinceFirstSeen    float64     `json:"since_first_seen"`
	Name              string      `json:"name"`
	ChassisID         string      `json:"chassis_id"`
	IP                string      `json:"ip"`
	OsVer             string      `json:"os_ver"`
	SinceLastDiscover float64     `json:"since_last_discover"`
}

// Gather device polling data
func PollingMetrics(ApiKey string) {
	path := "/api/v1/search/device?layers=3"
	url := netdiscoHost + path
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Error(err)
	}

	client := &http.Client{
		Timeout: time.Second * 10,
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", ApiKey)

	resp, err := client.Do(req)
	if err != nil {
		log.Error(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
	}

	data := Device{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Error(err)
	}

	// Generate Prometheus metrics for device polling
	for _, entry := range data {
		if len(entry.DNS) > 0 {
			NETDISCO_LAST_DISCOVER.WithLabelValues(entry.DNS).Set(entry.SinceLastDiscover)
			NETDISCO_LAST_ARPNIP.WithLabelValues(entry.DNS).Set(entry.SinceLastArpnip)
			NETDISCO_LAST_MACSUCK.WithLabelValues(entry.DNS).Set(entry.SinceLastMacsuck)
		}
	}

}
