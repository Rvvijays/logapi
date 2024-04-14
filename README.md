# LogAPI

This RESTful API allows you to search log lines based on specified criteria in a remote storage. The logs are stored in plain text format (.txt) in a structured folder hierarchy.

## Requirements

- Go (Golang)
- Any S3 Storage credentials

## Build and Run

1. Clone the repository:

   ```bash
   git clone https://github.com/Rvvijays/logapi.git
   cd logapi


2. Set up config:

    If you are using AWS S3, set your AWS credentials(endpoint, hostKey, region, secretKey), any S3 storage is compatible.
    ```bash
    {
        "endpoint": "https://s3.{storage_region}.wasabisys.com",
        "hostkey": "access_key",
        "region": "storage_region",
        "secretkey": "secret_key"
    }
    ```
    Make sure whatever bucket you are using is already created and have all the permission to create objects inside it with the provided credentials.
    - you can set the port | default: 8765
    - you can set interval(minutes) to write logs in local file. Default 5. If you set it to 0, automatic goroutine will not be started.
    - you can set interval(minutes) to upload local files to storage server. Default: 30, If you set it to 0, automatic goroutine will not be started.
    - you can enable/disable storage server upload by setting serverStorage to true/false respectively | default: false, after adding correct credentials make it true.

3. Buid and run the service:
    ```bash
    go build -o logapi .
    ./logapi

> The service should be running on `http://localhost:8765`:

## API Endpoints
### Add logs
- URL: `/api/logs/add`
- Method: POST
- Request Body:
    ```bash
    {
    "message": "log_message_here"
    }
- Sample Request
    ```bash
    curl -X POST http://localhost:8765/api/logs/add -d '{"message":"error: getting error in file opening."}'

### Search Logs
- URL: `/api/logs/search`
- Method: POST
- Request Body:
    ```bash
    {
        "searchKeyword":"your_search_keyword",
        "from":"yyyy-mm-dd",
        "to":"yyyy-mm-dd"
    }
- Sample request:
    ```bash
    curl -X POST http://localhost:8765/api/logs/search -d '{"searchKeyword":"error","from":"2024-04-01","to":"2024-04-01"}'



