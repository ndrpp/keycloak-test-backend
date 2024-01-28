# keycloak-test-backend

## Running locally

- Keycloak

Run with docker compose:

```
docker compose -f ./docker/docker-compose.yml up
```

Then, login on keycloak on localhost:8080 and create a realm, a client and a test user. Follow [this](https://www.keycloak.org/getting-started/getting-started-zip) if stuck.

- Go backend

Create an .env file in the project with the following variables taken from keycloak:
KEYCLOAK_REALM, KEYCLOAK_CLIENT_ID, KEYCLOAK_CLIENT_SECRET.
Then the server can be started using:

```
go build -o go-server
./go-server
```

- React frontend

Cd into the kc-front directory and run using:

```
bun run dev
```
