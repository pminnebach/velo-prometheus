package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-resty/resty"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	addr = flag.String("listen-address", ":8080", "The address to listen on for HTTP requests.")

	veloURI = "https://www.velo-antwerpen.be/availability_map/getJsonObject"
)

func init() {
	velo := NewVeloManager("antwerpen")
	prometheus.MustRegister(velo)
}

// Station holds all the data that a single station has.
// See Station.json for how the original looks like.
type Station struct {
	ID             string `json:"id"`
	District       string `json:"-"`
	Lon            string `json:"-"`
	Lat            string `json:"-"`
	Bikes          string `json:"bikes"`
	Slots          string `json:"slots"`
	Zip            string `json:"-"`
	Address        string `json:"-"`
	AddressNumber  string `json:"-"`
	NearbyStations string `json:"-"`
	Status         string `json:"-"`
	Name           string `json:"name"`
}

// VeloManager is what does a .
type VeloManager struct {
	City      string // Not realy necessary. Might only be usefull if another city uses the exact same infrastructure from ClearChanel.
	BikesDesc *prometheus.Desc
	SlotsDesc *prometheus.Desc
}

// GetStations calls the velo api and converts the json to whatever prometheus can handle.
func (v *VeloManager) GetStations() (bikesByStation map[string]int, slotsByStation map[string]int) {

	start := time.Now()

	// GET request
	resp, err := resty.R().Get(veloURI)
	if err != nil {
		fmt.Println("Error GET: ", err)
		return
	}

	var Stations []Station

	JSONObject := resp.Body()
	JSONBytes := []byte(JSONObject)

	// Unmarshal JSONBytes unto Stations.
	err = json.Unmarshal(JSONBytes, &Stations)
	if err != nil {
		log.Printf("Error unmashaling json: %v", err)
	}

	// Create 2 maps which will hold the values of bikes and slots.
	bikesByStation = make(map[string]int)
	slotsByStation = make(map[string]int)

	// Loop over all stations and pull out the available bikes,
	// slots and the station name.
	for idx := range Stations {
		bike := Stations[idx].Bikes
		slots := Stations[idx].Slots
		name := Stations[idx].Name

		// Convert json string to int.
		iBikes, err := strconv.Atoi(bike)
		if err != nil {
			log.Printf("Error strconv bike: %v", err)
		}

		iSlots, err := strconv.Atoi(slots)
		if err != nil {
			log.Printf("Error strconv slot: %v", err)
		}

		bikesByStation[name] = iBikes
		slotsByStation[name] = iSlots
	}

	end := time.Now()
	log.Printf("Duration: %v", end.Sub(start))

	return
}

// Describe simply sends the twp Descs in the struct to the channel.
func (v *VeloManager) Describe(ch chan<- *prometheus.Desc) {
	ch <- v.BikesDesc
	ch <- v.SlotsDesc
}

// Collect first triggers the GetStations. Then it creates the guages
// for each station on the fly based on the returned data.
func (v *VeloManager) Collect(ch chan<- prometheus.Metric) {

	bikesByStation, slotsByStation := v.GetStations()

	for station, bikesCount := range bikesByStation {
		ch <- prometheus.MustNewConstMetric(
			v.BikesDesc,
			prometheus.GaugeValue,
			float64(bikesCount),
			station,
		)
	}

	for station, slotsCount := range slotsByStation {
		ch <- prometheus.MustNewConstMetric(
			v.SlotsDesc,
			prometheus.GaugeValue,
			float64(slotsCount),
			station,
		)
	}
}

// NewVeloManager creates the two Descs BikesDesc and SlotsDesc.
func NewVeloManager(city string) *VeloManager {
	return &VeloManager{
		City: city,
		BikesDesc: prometheus.NewDesc(
			"velo_available_bikes",
			"Number of bikes available at a given station.",
			[]string{"station"},
			prometheus.Labels{"city": city},
		),
		SlotsDesc: prometheus.NewDesc(
			"velo_available_slots",
			"Number of free slots available at a given station.",
			[]string{"station"},
			prometheus.Labels{"city": city},
		),
	}
}

func main() {
	flag.Parse()

	log.Printf("Pushing metrics on port: %v", *addr)

	// Expose the registered metrics via HTTP.
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(*addr, nil))
}
