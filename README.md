# metricgen

metricgen is a simple Go app which allows the user to push a float value 
via /push. This value is consumed by the Prometheus Go library and exposed 
as a `test_input_number` gauge.

The purpose of this app is to allow a user to test their Prometheus scrape 
configuration.
