# Rates API

The Rates API provides an interface to fetch the latest and historical exchange rates for various currencies.

## Getting Started

These instructions will guide you through setting up and running the Rates API on your local machine for development and testing purposes.

### Prerequisites

- [Go](https://golang.org/doc/install) (Version 1.22)
- [Docker](https://docs.docker.com/get-docker/)

## Setup

The project uses Makefile to simplify the setup and running of the application. The [Makefile](Makefile) contains various commands to build, run, and clean up the project.

### Start Database Container:

The project requires a PostgreSQL database. A Docker Compose file is provided to run the database in a Docker container.

`make compose-up`

This command will set up a PostgreSQL server based on the configuration specified in database-docker-compose.yml.

### Application Configuration:

Ensure the application configuration, such as database connection settings, are correctly set in the application's configuration file : [config.yaml](config.yaml)

### Building the Project

To compile the project into a binary: `make build`
This will generate a binary named rates-api in the project directory.

### Running the Application

You can then run : `./rates-api --config-path config.yaml` to start the application.

### To start the application: `make run`

This command utilizes `go run` to start the application directly from the source code. Also it will start the application with the configuration file specified in the `config.yaml` file and the database container as specified.

### Cleaning Up

To remove generated files and stop the database container: `make clean`

## API Usage

- Fetch Latest Rates: [GET] /rates/latest
- Fetch Rates by Date: [GET] /rates/{date}
- Analyze Rates: [GET] /rates/analyze

More detailed API documentation is available at [open-api.spec.yaml](open-api.spec.yaml).
