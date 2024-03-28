package main

import (
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "postgres"
)

type ExchangeRate struct {
	XMLName xml.Name   `xml:"gesmes:Envelope"`
	Cube    CubeStruct `xml:"http://www.ecb.int/vocabulary/2002-08-01/eurofxref Cube"`
}

type CubeStruct struct {
	Time        string `xml:"time,attr"`
	CubeContent []Cube `xml:"Cube"`
}

type Cube struct {
	Currency string  `xml:"currency,attr"`
	Rate     float64 `xml:"rate,attr"`
	Time     string  `xml:"time,attr"`
}

var (
	db    *sql.DB
	mutex sync.RWMutex
)

func main() {
	initDB()
	loadExchangeRates()

	http.HandleFunc("/rates/latest", latestRatesHandler)
	http.HandleFunc("/rates/", specificDateRatesHandler)
	http.HandleFunc("/rates/analyze", analyzeRatesHandler)

	slog.Info("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func loadExchangeRates() {
	resp, err := http.Get("https://www.ecb.europa.eu/stats/eurofxref/eurofxref-hist-90d.xml")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var data ExchangeRate
	if err := xml.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Fatal(err)
	}

	for _, cube := range data.Cube.CubeContent {
		date, _ := time.Parse("2006-01-02", cube.Time)

		mutex.Lock()
		_, err := db.Exec("INSERT INTO exchange_rates (date, rates) VALUES ($1, $2) ON CONFLICT (date) DO NOTHING",
			date,
			fmt.Sprintf(`{"%s": %f}`, cube.Currency, cube.Rate))
		mutex.Unlock()

		if err != nil {
			log.Println("Error inserting into database:", err)
		}
	}
}

func latestRatesHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT rates FROM exchange_rates ORDER BY date DESC LIMIT 1")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var ratesJSON string
	for rows.Next() {
		err := rows.Scan(&ratesJSON)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, ratesJSON)
}

func specificDateRatesHandler(w http.ResponseWriter, r *http.Request) {
	date := r.URL.Path[len("/rates/"):]
	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		http.Error(w, "Invalid date format", http.StatusBadRequest)
		return
	}

	var ratesJSON string
	err = db.QueryRow("SELECT rates FROM exchange_rates WHERE date = $1", t).Scan(&ratesJSON)
	switch {
	case err == sql.ErrNoRows:
		http.Error(w, "Exchange rates not available for the specified date", http.StatusNotFound)
		return
	case err != nil:
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, ratesJSON)
}

func analyzeRatesHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT rates FROM exchange_rates")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type rateData struct {
		Min float64
		Max float64
		Avg float64
	}

	analyzeData := make(map[string]*rateData)

	for rows.Next() {
		var ratesJSON string
		if err := rows.Scan(&ratesJSON); err != nil {
			log.Println("Error scanning row:", err)
			continue
		}

		var rates map[string]float64
		if err := json.Unmarshal([]byte(ratesJSON), &rates); err != nil {
			log.Println("Error unmarshalling JSON:", err)
			continue
		}

		for currency, rate := range rates {
			data, ok := analyzeData[currency]
			if !ok {
				data = &rateData{}
				analyzeData[currency] = data
			}

			if data.Min == 0 || rate < data.Min {
				data.Min = rate
			}
			if rate > data.Max {
				data.Max = rate
			}
			data.Avg += rate
		}
	}

	response := make(map[string]map[string]float64)
	for currency, data := range analyzeData {
		response[currency] = map[string]float64{
			"min": data.Min,
			"max": data.Max,
			"avg": data.Avg / float64(len(analyzeData)),
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
