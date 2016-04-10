# Particle Cloud API

[![Circle CI](https://circleci.com/gh/sepal/particle.svg?style=svg)](https://circleci.com/gh/sepal/particle)

[Go](https://golang.org/) API client for the [Particle Cloud](https://www.particle.io/).

**It is currently WIP.**

## Todo

- ~~Device Info~~
- ~~Device Variables~~
- ~~Device Functions~~
- ~~Events~~
- Refactor/Improve error handling
- Special Events

Further functionality is not planned for now, but might be 
added later. Pull Requests are of course welcomed :-)


## Installation

```bash
go get github.com/sepal/particle
```

## Running the examples

1. Programm your core/photon/electron with the firmware in 
`examples/device/firmware/test.ino`

2. Go into one of the examples and either run `go install` to install 
the example to your $GOPATH/bin path or execute the example directly 
like this:

```
cd examples/device/list_devices
go run list_devices.go -t your_token
```