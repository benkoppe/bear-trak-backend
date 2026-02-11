# bear-trak-backend

This is the backend for BearTrak!

It consists of a Go server and [OpenTripPlanner](https://www.opentripplanner.org/) (otp) instance. Each is individually deployed with [Komodo](https://www.komo.do/) on mostly Oracle infrastructure.

## How to deploy

Each piece of the backend has its own Dockerfile and is built into its own image. See each subdirectory for more information.

**Pushes to the `main` branch trigger a GitHub action that builds & pushes both images. Komodo pulls these images and updates automatically.**
