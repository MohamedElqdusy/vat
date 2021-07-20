# VAT Validation

## System Requirements
- Docker

## Usage 
#### Run
    > docker build -t vat .
    > docker run -p 7888:7888 vat
    
check `usage.http` for the api usage, just click over them.

###  Integration test
    > go test ./...
 ---
## Notes:
- all the written tests are integrations so you need to be connected to the internet.
