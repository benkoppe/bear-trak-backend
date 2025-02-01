# bear-trak-backend

This is BearTrak's backend.

It consists of a Go server and [OpenTripPlanner](https://www.opentripplanner.org/) (otp) instance.

## How to deploy

Pushes to the `main` branch trigger a GitHub action that builds & pushes both images, then updates the deployment with `docker stack`.

- When first setting up the server side, follow these steps:
  - Create a new user `deploy` with ssh access granted to a public key in the `.ssh/authorized_keys` file.
  - Before the allowed public key, add `command="docker system dial-stdio"` to `.ssh/authorized_keys` to restrict command usage.
  - Set the repository secret `DEPLOY_SSH_PRIVATE_KEY` to the corresponding private key.
  - You're done! A `docker stack deploy` command will be used to push production to the server.
  - To bypass the GitHub action, you can set up your own Docker context and run `docker stack deploy -c ./docker-stack.yml bear-trak`.
- Now that a server has been set up, the GitHub action will automatically deploy using `docker-stack.yml` :)
