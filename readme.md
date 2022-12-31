# Build your own Cosmos Blockchain
https://tutorials.cosmos.network/hands-on-exercise/1-ignite-cli/3-stored-game.html#

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

```
# (Re)generate protobuf
docker exec -it checkers ignite generate proto-go
```

```
# All unit test
docker exec -it checkers go test -v ./...
# Specific dir
docker exec -it checkers go test -v go test github.com/dan-l/checkers/x/checkers/types
# Clean cache
docker exec -it checkers go clean -testcache
```

### Faucet
http://localhost:4500/


## Proto

- `query` : reading state
- `tx`: update state
- `genesis`: genesis state

## State

### System info 
- See struct at `x/checkers/types/system_info.pb.go`
- Store system wide info such as the next id for new game (which also tracks number of game history) 
- `docker exec -it checkers checkersd query checkers show-system-info -o json`

### Stored Game
- See struct at `x/checkers/types/stored_game.pb.go`
- Information saved about the game state such as players and the board state
- Map keyed by index to query the game state by index
- `docker exec -it checkers  checkersd query checkers list-stored-game -o json`