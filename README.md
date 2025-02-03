# Backend for Image Processing

## Description

Written purely in Go, this service helps to process images from payload, so that they can be stored efficiently in database and can be accessed in minimum latency. It downloads the image, calculates its perimeter and stores the result in database. Uses Read-Write mutex and concurrency powers of Golang, can run locally or with docker.

## Assumptions

- Only jobID and status are necessary fields, other fields like storeID/failedID, imageURL or results can be updated as processed.
- Images are of 480p, which means approximately width and height are 800:400.
- Master Store data can be stored and used locally for quick look-up of storeID for now.

## Work Environment

- Ubuntu distro over WSL (on Windows 11)
- VSCode text-editor
- Gin web framework, MongoDB for database

## Installing & Testing

### Installation

1. **Clone Repository**:
```sh
git clone github.com/prkshayush/img-processing
cd img-processing
```
2. set-up .env file variables:
MONGODB_URI - Connection string, starting with 'mongodb+...'
MONGODB_DB - Database name
MONGODB_COLLECTION - Collection Name of MongoDB.
MASTER_STORE_PATH - file path where Master Store data is saved.

3. Run in docker container:
```sh
docker build -t img-processing .
docker run -d -p 8080:8080 --name img-processing \
  --env-file .env \
  -v $(pwd)/data:/data \
  img-processing
```

### Testing

1. submit job: 
```sh
curl -X POST http://localhost:8080/api/submit -H "Content-Type: application/json" -d '{
  "count": 2,
  "visits": [
    {
      "store_id": "S00339218",
      "image_url": [
        "https://www.gstatic.com/webp/gallery/2.jpg",
        "https://www.gstatic.com/webp/gallery/3.jpg"
      ],
      "visit_time": "time of store visit"
    },
    {
      "store_id": "S01408764",
      "image_url": [
        "https://www.gstatic.com/webp/gallery/3.jpg"
      ],
      "visit_time": "time of store visit"
    }
  ]
}'
```

2. Check status of job using jobID from previous response:
```sh
curl -X GET "http://localhost:8080/api/status?jobID=<job_id>"
```
