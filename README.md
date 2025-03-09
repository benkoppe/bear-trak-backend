# bear-trak-backend

This is the backend for BearTrak!

It consists of a Go server and [OpenTripPlanner](https://www.opentripplanner.org/) (otp) instance. Each is individually deployed with [Coolify](https://www.coolify.io/).

## Todo

### Planned features

- [ ] gym capacity history
- [ ] dining menu favorites
- [ ] libraries

### Backlog

- [ ] write `go-server` tests & integrate into `Pipeline`.
- [ ] centralize shared data like external transit URLs.
- [ ] figure out how to detect when the `otp` static GTFS data becomes invalid.
- [ ] fix incorrect GTFS-RT vehicle positions being served to `otp` directly from availtec.

## How to deploy

Each piece of the backend has its own Dockerfile and is built into its own image. See each subdirectory for more information.

**Pushes to the `main` branch trigger a GitHub action that builds & pushes both images, then updates the deployment with a call to Coolify.**

### New server setup

- Run Coolify's server connection setup on your server
- Create two `Docker Compose` resources -- one for the go-server, one for the otp instance. Set base folders accordingly.
  - On Coolify, set the Docker build command for both resources to `docker compose pull` to ensure that images are pulled, not built.
- In the resource for `go-server`, set the following env variable:

| Name                         | Description                  | Example value |
| ---------------------------- | ---------------------------  | ------------- |
| POSTGRES_PASSWORD            | Password for the Postgres db | 123456...     |

- Then, in the GitHub repository, set the following secrets:
  
| Name                         | Description                  | Example value |
| ---------------------------- | ---------------------------  | ------------- |
| COOLIFY_API_TOKEN            | Write/Deploy API token for Coolify | uqpDqZK6...     |
| COOLIFY_URL            | URL for your Coolify instance | coolify.yourdomain.com     |

  - As well as these variables:

| Name                         | Description                  | Example value |
| ---------------------------- | ---------------------------  | ------------- |
| GO_PRODUCTION_RESOURCE_ID            | `go-server` production resource id in Coolify | Qv2BXOMteC3h...     |
| OTP_PRODUCTION_RESOURCE_ID            | `otp` production resource id in Coolify | 2tuA1wlYex...     |
| GO_STAGING_RESOURCE_ID            | `go-server` staging resource id in Coolify (optional) | 3sLCDmy6MAC...     |
| GO_STAGING_RESOURCE_ID            | `otp` staging resource id in Coolify (optional) | LODIKoD9x...     |

- You're done! The `Pipeline` workflow will use Coolify to deploy images on your server every time you push. :)
  - NOTE: The workflow will fix the commit hash in Coolify, so if manually updating, make sure you reset to latest!
