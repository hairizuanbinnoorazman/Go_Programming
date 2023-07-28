# Basic Using Datastore

This is a basic application that tries to using Google Datastore locally.


## Quickstart

The following command starts datastore that we can then use Golang code to test against it

```bash
docker run -p 8081:8081 google/cloud-sdk:437.0.1 gcloud beta emulators datastore start --project=test --host-port=0.0.0.0:8081
```

To start application

```bash
DATASTORE_PROJECT_ID=test DATASTORE_EMULATOR_HOST=localhost:8081 go run main.go
```
