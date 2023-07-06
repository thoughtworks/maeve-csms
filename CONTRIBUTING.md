# How to contribute to MaEVe

Thank you for taking the time to read this guide. We are happy you wish to contribute to this project :partying_face: !

We value all contributions, and we hope the following sections will give you clarity on how you can contribute.

## Code of Conduct

This project adheres to the Contributor Covenant code of conduct. By participating, you are expected to uphold this
code.

Our code of conduct is available
at [![Contributor Covenant](https://img.shields.io/badge/Contributor%20Covenant-2.1-4baaaa.svg)](./CODE_OF_CONDUCT.md)

## Requirements

Please take some time to familiarise yourself with the components in MaEVe. These are in
the [README's documentation section](./README.md/#documentation).

You will need to have Go 1.20 installed in your development machine. Instructions are available on
the [Go website](https://go.dev/doc/install)

MaEVe components can be run through a [docker-compose file](./docker-compose.yml), hence you will require Docker on your
machine. We mostly use MacOS and run the docker daemon using [Colima](https://github.com/abiosoft/colima).

Additionally, you will require the [compose plugin](https://docs.docker.com/compose/install) to run the containers
listed in the [docker-compose file](./docker-compose.yml).

### What dependencies are provisioned through docker compose?

Currently, the [docker-compose file](./docker-compose.yml) provisions the following
in addition to the [Manager](./docs/manager.md) and [Gateway](./docs/gateway.md) components:

- #### MQTT v5 broker
  This is used to decouple the stateful connections required by OCPP from the stateless message handlers. The system
  currently connects to the MQTT broker anonymously.

- #### Redis
  The system uses [Redis](https://redis.io/) as a storage engine for transaction details. The manager connects to the
  store anonymously.

- #### Load balancer
  A load balancer implemented through [envoyproxy](https://www.envoyproxy.io/) to load balance calls coming from the
  charge stations and allowing the system to perform better under high demand.

## Build Go projects locally

Whilst the components can be run through Docker, if you wish to run locally without containerising them, you can run the
following command:

```shell
$ go build ./...
```

_You will need to run this from the correct sub-directory (e.g. ./gateway or ./manager)_

## Testing

Each Go project contains a set of tests that are when a commit reaches the remote repository.
These tests also run as part of our [git hooks](#install-the-hooks) before committing to the local repository.

If you wish to run the tests manually, you can run the following:

```shell
$ export TESTCONTAINERS_DOCKER_SOCKET_OVERRIDE=/var/run/docker.sock
$ export DOCKER_HOST=$(docker context inspect -f '{{ .Endpoints.docker.Host }}')
$ go test ./...
```

_Please ensure the docker daemon is running on your machine before running these tests._

## How to report a bug and track issues

Bugs are reported and tracked via [GitHub issues](https://github.com/twlabs/maeve-csms/issues). Please use this platform
to report any issues to our community.

When reporting a bug, please provide the following information about your environment setup:

- Operating System (OS) and version
- Code version you are using
- Steps to reproduce the issue and any additional guidance

Our team will periodically review and prioritize reported issues for resolution.

For support, please contact us at [maeve-team@thoughtworks.com](mailto:maeve-team@thoughtworks.com)

## How to submit changes

Please follow these steps to submit your changes:

1. Fork the repository by clicking the "Fork" button on the [repository's page](https://github.com/twlabs/maeve-csms).

2. Work on your own repository fork following our guidelines listed on this page until you are satisfied with your
   work (e.g. a feature has been completed, or a bug is now fixed)
3. Create a pull request on this repository and provide a descriptive title with a detailed description of your changes.
   Additionally, please make sure that you reference the issue, if present, that your contribution is solving.

Our team will be happy to review your changes and work together to make them available to the rest of the community.

## Coding style guidelines

Please make sure you have configured the following hooks on your development environment. They ensure that, in addition
to
scanning the code for vulnerabilities and code smells, all contributors adopt the same code style.

### Install the hooks

You will first need to install [the pre-commit framework](https://pre-commit.com/#install) in order to configure the git
hooks.

Then, from the root of the repo, you can run the following:

```shell
$ go install honnef.co/go/tools/cmd/staticcheck@latest
$ go install golang.org/x/tools/cmd/goimports@latest
$ go install github.com/securego/gosec/v2/cmd/gosec@latest
$ pre-commit install && pre-commit install --hook-type commit-msg
```

## Commit messages

Commit messages should follow the [Conventional Commits specification](https://www.conventionalcommits.org/en/v1.0.0/).

Additionally, please ensure that the _body_ of your commit message explicitly and succinctly communicates to other
contributors what your change is about.
This will us keep commits small and centered around something specific we are working on.

For example, you can check
out [this commit](https://github.com/twlabs/maeve-csms/commit/2c64552a689f185728566c841bfa7609469015f5) as a reference
to what a good commit looks like.

_Please note that all verbs should be at the present tense e.g. Add / ~~Adding~~ / ~~Added~~ ._

## Contributors

_Come back to this later, we are currently working on it._