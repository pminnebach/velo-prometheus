# velo-prometheus

The Velo station information is pushed as a [json](https://www.velo-antwerpen.be/availability_map/getJsonObject) feed. I created this, because my [initial project](https://github.com/pminnebach/Velo) to push these metrics to InfluxDB was kinda buggy and unstable.

And because statistics are awesome.

## Usage

### local

````
$ go build -o velo
````

````
$ ./velo --help

Usage of ./velo:
  -listen-address string
        The address to listen on for HTTP requests. (default ":8080")
````

### Docker

````
$ docker build -t velo .
$ docker run -d -p 8080:8080 velo
````

In your prometheus configuration, add the following to `scrape_configs:`

````
  - job_name: 'velo'
    scrape_interval: 10s
    static_configs:
      - targets: ['<url/ip>:8080']
````

## Todo

- [ ] Make Dockerfile GOARCH independent.
- [ ] Remove Resty dependency.
- [ ] Add fancy CLI flags. (Cobra?)
- [ ] Add error handeling where necessary.

## Next

To build upon this i want something that can use those metrics to suggest me which station i should go to for getting or putting back a bike.
This based on my current location (home or work).

Or maybe in the future by placing a marker on a map and it figures out which stations is the closest based on lat/lon, and suggests nearby stations if the chosen one is full or empty.

## Contributing

Send PR.

## Disclaimer

I'm not a professional developer.
