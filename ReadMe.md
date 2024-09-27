# Run the application locally using Docker Compose
docker-compose up --build

# Run app on Portainer - requires a personal Docker registry
cd deploy <!-- Go to the deploy directory -->
cp .env.sample .env <!-- Copy the sample .env file and fill data -->
make build_dev <!-- Build the development version of the application -->

- go to `Stacks` (in Portainer)
- add stack
- copy content of `docker-compose.yml`
- in `Environment variables` click on `Advanced mode` and copy the `volumes/app/.env` content
- click on `Deploy the stack`
