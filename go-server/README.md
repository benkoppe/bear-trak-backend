# BearTrak's `go-server`

This Go server provides BearTrak with most of its data!

Start it up with:

```bash
go run main.go
```

This will serve both an OpenAPI-specified REST API (described below), and some static files on the subpath `/static`.

It will also start a (currently unused) hourly timer to run tasks.

To build and run the Docker container:

```
docker compose up --build
```

## OpenAPI Code Generation

The OpenAPI spec file [openapi.yaml](/go-server/openapi.yaml) is the source of truth for this backend. Both the client and server sides are generated from this configuration file.

`go-server` uses [ogen](https://github.com/ogen-go/ogen) to generate the `api` module, which provides contracts for the backend to fulfill.

To regenerate the `api` module, see instructions in `generate.go`.
