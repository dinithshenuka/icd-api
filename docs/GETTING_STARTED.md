# Getting Started

Follow these instructions to set up the ICD Codes API in your local environment.

## Prerequisites
- Go 1.22 or higher
- Git

## Installation
1. Clone the repository:
   ```bash
   git clone https://github.com/iammrdp/icd-api.git
   cd icd-api
   ```
2. Install dependencies:
   ```bash
   go mod download
   ```

## Running the API
To start the service, execute the following command from the root directory:
```bash
go run cmd/api/main.go
```

The server will be available at `http://localhost:8080`.

## Testing the Endpoints
You can verify the installation by retrieving a sample of ICD-10 codes:
```bash
curl http://localhost:8080/icd10
```
