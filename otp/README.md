# BearTrak's `otp` instance

This is a deployment of OpenTripPlanner, set to use:
- a cropped PBF file for Ithaca, downloaded into this repo in December 2024
- a static GTFS zip for TCAT bus & schedule data, downloaded into the Docker container at build-time
- a `router-config.json` file that allows it to source realtime GTFS data for live TCAT predictions

Currently, `otp` is deployed to a separate domain, and used by clients separately from `go-server`. The two deployments don't interact.

The two are kept separate mostly because `otp` is only used for routing predictions. Though `go-server` could wrap around `otp`, with its outputs worked into the OpenAPI schema, I didn't want to add another layer of abstraction unnecessarily. In the future, if that design changes, this instance may become cut off from outside requests.

See the Dockerfile to understand how the instance is configured.

## Details

Feed appears to have ID: 1

### Realtime GTFS data

Three realtime types:

- Alert
- VehiclePosition
- TripUpdate

<https://realtimetcatbus.availtec.com/InfoPoint/GTFS-Realtime.ashx?&Type={TYPE}>
