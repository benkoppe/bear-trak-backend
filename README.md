# bear-trak-backend

This is the backend for BearTrak.

It consists of a Go server and [OpenTripPlanner](https://www.opentripplanner.org/) (otp) instance.

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

Each piece of the backend is containerized - latest images can be seen in the `Packages` section on the right. Images are pulled and reverse proxied with [traefik](https://github.com/traefik/traefik) in accordance to `docker-stack.yml`.

**Pushes to the `main` branch trigger a GitHub action that builds & pushes both images, then updates the deployment with `docker stack`.**

### New server setup

- On your own machine:
  - Create a new SSH key pair with `ssh-keygen -t ed25519 -C "deploy@{serverIP}"`.
- On the server:
  - Install Docker.
  - Create a new user `deploy`. Switch users with `su - deploy`.
  - Modify `.ssh/authorized_keys`:
    - On a single line, first type `command="docker system dial-stdio"` to restrict key access.
    - On the same line, paste the contents of the `.pub` key file from your newly created SSH key pair.
  - Run `docker swarm init`.
- In this repository:
  - Set the repository secret `DEPLOY_SSH_PRIVATE_KEY` to the private key file from your newly created SSH key pair.
  - Also set secret `CLOUDFLARE_DNS_API_TOKEN` to an API key with permissions `DNS:Edit`.
  - You're done! The `Pipeline` workflow will use `docker stack deploy` to deploy images on your server.

Now that a server has been set up, the GitHub action will keep the server up-to-date with every push.  :)

### Bypass the GitHub action
- Set up a docker context with `docker context create {name} --docker "host=ssh://{user}@{serverIP}`.
- Use the context with `docker context use {name}`.
  - NOTE: With this context active, your docker commands will be run on the server.
- Make sure you've set a local `CLOUDFLARE_DNS_API_TOKEN`.
- Run `docker stack deploy -c ./docker-stack.yml bear-trak`.
