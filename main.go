package main

import (
        "fmt"
        "io"
        "net/http"
        "strconv"

        "github.com/prometheus/client_golang/prometheus"
        "github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
        testInput = prometheus.NewGauge(prometheus.GaugeOpts{
                Name: "test_input_number",
                Help: "Latest input to /push.",
        })
)

func init() {
        prometheus.MustRegister(testInput)
}

func push(w http.ResponseWriter, r *http.Request) {
        hasInput := r.URL.Query().Has("num")
        inputString := r.URL.Query().Get("num")
        inputFloat, err := strconv.ParseFloat(inputString, 8)

        fmt.Printf("[%s] - got / request. push(%t)=%s\n",
                r.RemoteAddr, hasInput, inputString)

        if !hasInput {
                io.WriteString(
                        w, 
                        "Use 'num=____' query parameter to push a float\n",
                )
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

        http.HandleFunc(
                "/", 
                func(w http.ResponseWriter, _ *http.Request) {
                        io.WriteString(
                                w, 
                                "metricgen allows a user to input a float value" +
                                " to /push which is consumed by Prometheus.\n",
                        )
                },
        )
        http.Handle("/metrics", promhttp.Handler())
        http.HandleFunc("/push", push)

        http.ListenAndServe(":2112", nil)
}