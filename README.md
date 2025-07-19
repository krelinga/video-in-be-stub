# video-in-be-stub
Stub implementation of the video-in-be

## Running with Docker

### Build the Docker image
```bash
docker build -t video-in-be-stub .
```

### Run the container
```bash
docker run --rm -p 8080:8080 video-in-be-stub
```

The service will be available at `http://localhost:8080`

## Development

### Build locally
```bash
go build -o video-in-be-stub main.go
```

### Run locally
```bash
./video-in-be-stub
```
