# Build your own Cosmos Blockchain
**checkers** is a blockchain built using Cosmos SDK and Tendermint and created with [Ignite CLI](https://ignite.com/cli).

## Get started

```
ignite chain serve
```

`serve` command installs dependencies, builds, initializes, and starts your blockchain in development.

### Configure

Your blockchain in development can be configured with `config.yml`. To learn more, see the [Ignite CLI docs](https://docs.ignite.com).

### Web Frontend

Ignite CLI has scaffolded a Vue.js-based web app in the `vue` directory. Run the following commands to install dependencies and start the app:

```
cd vue
npm install
npm run serve
```

The frontend app is built using the `@starport/vue` and `@starport/vuex` packages. For details, see the [monorepo for Ignite front-end development](https://github.com/ignite-hq/web).

UI at http://localhost:3000/.

### Install
To install the latest version of your blockchain node's binary, execute the following command on your machine:

```
curl https://get.ignite.com/dan-l/checkers@latest! | sudo bash
```
`dan-l/checkers` should match the `username` and `repo_name` of the Github repository to which the source code was pushed. Learn more about [the install process](https://github.com/allinbits/starport-installer).

### Using docker

```
# Build image
docker build -f Dockerfile . -t checkers_i
# Create container
docker create --name checkers -i -v $(pwd):/checkers -w /checkers -p 1317:1317 -p 3000:3000 -p 4500:4500 -p 5000:5000 -p 26657:26657 checkers_i
# Start container
docker start checkers
# Scaffold
docker run --rm -it -v $(pwd):/checkers -w /checkers checkers_i ignite scaffold chain github.com/dan-l/checkers
```

```
# Start chain
docker exec -it checkers ignite chain serve
# Print status
docker exec -it checkers bash -c "checkersd status 2>&1 | jq"
```

```
# Help
docker exec -it checkers checkersd --help
docker exec -it checkers checkersd status --help
docker exec -it checkers checkersd query --help  
```

```
# Install
docker exec -it checkers bash -c "cd vue && npm install"
# Build and run UI
docker exec -it checkers bash -c "cd vue && npm run dev -- --host"
```

### Faucet
http://localhost:4500/

