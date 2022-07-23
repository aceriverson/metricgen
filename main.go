package main

import (
        "fmt"
        "io"
        "net/http"
        "strconv"
        // "time"

        "github.com/prometheus/client_golang/prometheus"
        // "github.com/prometheus/client_golang/prometheus/promauto"
        "github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
        testInput = prometheus.NewGauge(prometheus.GaugeOpts{
                Name: "test_input_number",
                Help: "Latest input to /push.",
        })
        // hdFailures = prometheus.NewCounterVec(
        //         prometheus.CounterOpts{
        //                 Name: "hd_errors_total",
        //                 Help: "Number of hard-disk errors.",
        //         },
        //         []string{"device"},
        // )
)

func init() {
        prometheus.MustRegister(testInput)
}

func push(w http.ResponseWriter, r *http.Request) {
        hasInput := r.URL.Query().Has("num")
        inputString := r.URL.Query().Get("num")
        inputFloat, err := strconv.ParseFloat(inputString, 8)

        fmt.Printf("got / request. push(%t)=%s\n",
                hasInput, inputString)

        if !hasInput {
                io.WriteString(w, "Use 'num=____' query parameter to push a string\n")
        } else if err != nil {
                io.WriteString(w, "'num' query must be a float\n")
        } else {
                testInput.Set(inputFloat)
                io.WriteString(w, fmt.Sprintf("Pushed %f\n", inputFloat))
        }
}

func main() {
        testInput.Set(0)

        fmt.Println("metricgen Started on Port 2112")

        http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {io.WriteString(w, "metricgen allows a user to input a float value to /push which is consumed by Prometheus.\n")})
        http.Handle("/metrics", promhttp.Handler())
        http.HandleFunc("/push", push)

        http.ListenAndServe(":2112", nil)
}