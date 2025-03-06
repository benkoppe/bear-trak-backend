# BearTrak's `go-server`

This Go server provides BearTrak with most of its data!

Start it up with:

```bash
go run main.go
```

This will serve both an OpenAPI-specified REST API (described below), and some static files on the subpath `/static`.

The go server requires a connection to a postgres database. See `compose.yaml` for required environment variables. Schema modifications are made in `/migrations`, and keep up-to-date in the database automatically. 

It will also start a five-minute timer to run tasks, currently used for logging gym capacity updates.

The Docker container contains two images: `go-server` and `postgres`. They can both be started with:

```
docker compose up --build
```

This `compose.yaml` is also used by the Coolify instance, with the password set using an environment variable.

## OpenAPI Code Generation

The OpenAPI spec file [openapi.yaml](/go-server/openapi.yaml) is the source of truth for this backend. Both the client and server sides are generated from this configuration file.

`go-server` uses [ogen](https://github.com/ogen-go/ogen) to generate the `api` module, which provides contracts for the backend to fulfill.

To regenerate the `api` module, see instructions in `generate.go`.
