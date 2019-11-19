# statuscake-exporter

Prometheus StatusCake exporter

Exporter [consumes data from StatusCake API](https://www.statuscake.com/api/Period%20Data/Get%20Period%20Data.md) using the [official lib](https://godoc.org/github.com/DreamItGetIT/statuscake) exposing it to Prometheus on port X.

Supported metrics:

`statuscake_test_up`: Current status at last test
`statuscake_test_uptime`: 7 Day Uptime

## BUILD

`make build`

The binary will be created on `./bin` dir.

## OPTIONS

`-web.listen-address`: Exporter HTTP listening port
`-web.telemetry-path`: Exporter metrics path
`-stk.username`: StatusCake Username
`-stk.apikey`: StatusCake APIKEY
`-stk.tags`: Test's Tags to limit the scrape
`-stk.interval`: Interval in seconds to gather metrics from API
`-stk.enable-tests`: Enable Test module collector
`-stk.enable-ssl`: Enable SSL module collector
`-stk.ssl-flags`: SSL Flags to be enabled as metrics
`-v`: Show the build version

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

```bash
docker run -it mtulio/statuscake-exporter:v0.3.1 \
  -stk.username $STATUSCAKE_USER \
  -stk.apikey $STATUSCAKE_APIKEY
```

OR, to build your own Docker image:

`make build`

The binary will be created on `./bin` dir.

`docker build  .`

This will create a docker image.

Set STATUSCAKE_APIKEY and STATUSCAKE_USER variables for Docker image, and  `docker run $image-id-here`

## CONTRIBUTOR

You can contribute with three ways: using, testing and developing.

* Fork me
* Open an PR with enhancements, bugfixes, etc
* Open an issue

You are welcome. =)
