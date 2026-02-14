# bear-trak-backend

This is the backend monorepo for BearTrak!

It consists of:
- `go-server`: The primary backend. Almost all requests go to this server.
- `otp`: Configured instance of [OpenTripPlanner](https://www.opentripplanner.org/). Serves directions requests only.
- `transitclock`: Extended [TheTransitClock](https://thetransitclock.github.io/), designed to consume responses from `go-server`.
  - Required to support `otp` for schools without available GTFS-Realtime files (umich).

Each service is packaged with nix, allowing uniform commands to build development shells and output builds & Docker containers:
```bash
nix develop ./<folder> --no-pure-eval # enter dev shell
nix build ./<folder>#<output> # build named output
```

Everything is individually deployed with [Komodo](https://www.komo.do/) on mostly Oracle infrastructure.

## How to deploy

Each piece of the backend has its own nix flake and is built into its own image. See each subdirectory for more information.

**Pushes to the `main` branch trigger a GitHub action that builds & pushes both images. Komodo pulls these images and updates automatically.**
