# statuscake-exporter

Prometheus StatusCake exporter

Exporter [consumes data from StatusCake API](https://www.statuscake.com/api/Period%20Data/Get%20Period%20Data.md) using the [official lib](https://godoc.org/github.com/DreamItGetIT/statuscake) exposing it to Prometheus on port X.

Supported metrics:

`statuscake_test_up`: Current status at last test
`statuscake_test_uptime`: 7 Day Uptime

## BUILD

`make build`

The binary will be created on `./bin` dir.

## USAGE

* Show metrics from all StatusCake Tests

`./bin/statuscake-exporter -stk.username my_stk_user -stk.apikey my_stk_apikey`

```
# HELP statuscake_test_up Status Cake test Status
# TYPE statuscake_test_up gauge
statuscake_test_up{name="MyApp01_-_api"} 1
statuscake_test_up{name="MyApp02_-_front"} 1
# HELP statuscake_test_uptime Status Cake test Uptime from the last 7 day
# TYPE statuscake_test_uptime gauge
statuscake_test_uptime{name="MyApp01_-_api"} 100
statuscake_test_uptime{name="MyApp02_-_front"} 100
```

* Show metrics filtering by Tags from StatusCake Tests

`./bin/statuscake-exporter -stk.username my_stk_user -stk.apikey my_stk_apikey -stk.tags "MyApp01"`

```
# HELP statuscake_test_up Status Cake test Status
# TYPE statuscake_test_up gauge
statuscake_test_up{name="MyApp01_-_api"} 1
# HELP statuscake_test_uptime Status Cake test Uptime from the last 7 day
# TYPE statuscake_test_uptime gauge
statuscake_test_uptime{name="MyApp01_-_api"} 100
```

## USAGE DOCKER

> TODO

## CONTRIBUTOR

* Fork me
* Open an PR with enhancements, bugfixes, etc
* Open an issue

You are welcome. =)
