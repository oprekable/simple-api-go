# simple-api-go

simple-api-go

## How To Run
### Requirements
- Golang at least version `1.21`
- Docker
- `Make` command line tool

### Docker Up and Running
- To fire up mysql server do run command `make docker-up`
- To stop mysql server do run command `make docker-down`
- To reset/destroy mysql server data do run command `make docker-clean-data` (please do stop mysql server first)

### Migrate DB and seed
- Install `sql-migrate` with run command `make migrate-dev-install`
- Migrate and Seed DB data with run command `make migrate-dev-up`

### Run Application and Test API
- To run the application please make sure 
  - Docker already up and running
  - run command `make docker-up`
  - run command `make migrate-dev-up`
- Run the application with run command `make run-app`
- HTTP server should be running in port `3000`
- API address available
  - For no filter `http://localhost:3000/api` do test with run command `make test-api-no-filter`
  - For brand filter `http://localhost:3000/api?brand=Honda` do test with run command `make test-api-filter-brand`
  - For type filter `http://localhost:3000/api?type=Beat` do test with run command `make test-api-filter-type`
  - For transmission filter `http://localhost:3000/api?transmission=Manual` do test with run command `make test-api-filter-transmission`
  - For all fields filter `http://localhost:3000/api?brand=Honda&type=Beat&transmission=Automatic` do test with run command `make test-api-filter-all-fields`

### Run Everything in GitPod
- Go to https://gitpod.io/new/?autostart=true#https://github.com/oprekable/simple-api-go and login with your github account
- Wait until new tab automatically opened URL with address path such https://xxxx.xxx.gitpod.io/
- Access the API via https://xxxx.xxx.gitpod.io/api