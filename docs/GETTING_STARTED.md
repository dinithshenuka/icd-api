# Getting Started

Follow these instructions to set up the ICD Codes API in your local environment.

## Prerequisites
- Go 1.22 or higher
- Git
- Make

## Installation
1. Clone the repository:
   ```bash
   git clone https://github.com/dinithshenuka/icd-code-api.git
   cd icd-api
   ```

2. Initialize the Database and Import Data:
   ```bash
   make migrate-up
   make import
   ```

## Development
This project uses a `Makefile` to simplify common tasks.

### Start the Development Server
To start the server with live-reloading (Air):
```bash
make dev
```
The server will be available at `http://localhost:8080`.

### Generate Code
To regenerate Go code from the OpenAPI specification after making changes to `openapi/openapi.yaml`:
```bash
make generate
```

### Running Tests
To run the project's test suite:
```bash
make test
```

### Build for Production
To create a production binary:
```bash
make build
```

## Testing the Endpoints
You can verify the installation by searching for a code:
```bash
curl "http://localhost:8080/v1/search?q=heart"
```

Interactive documentation is available at `http://localhost:8080/docs`.
