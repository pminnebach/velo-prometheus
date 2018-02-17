# velo-prometheus

The Velo station information is pushed as a [json](https://www.velo-antwerpen.be/availability_map/getJsonObject) feed. I created this, because my [initial project](https://github.com/pminnebach/Velo) to push these metrics to InfluxDB was kinda buggy and unstable.

And because statistics are awesome.

## Usage

Run local

````
$ go build -o velo
````

````
$ ./velo --help

Usage of ./velo:
  -listen-address string
        The address to listen on for HTTP requests. (default ":8080")
````

Or as a docker container. 
I created this to run on a raspberry pi. So first change the GOARCH in the Dockerfile to the desired platform.

Port 8080 is exposed by default. To change this, add the `-listen-address [port]` parameter

````
$ docker build -t velo .
$ docker run -d -p 8080:8080 velo
````

## Todo

- [ ] Make Dockerfile GOARCH independent.
- [ ] Remove Resty dependency.
- [ ] Add fancy CLI flags. (Cobra?)
- [ ] Add graceful shutdown.
- [ ] Add error handeling where necessary.

## Contributing

Send PR.
