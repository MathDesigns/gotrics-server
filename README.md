# GoTrics Server

The **GoTrics Server** is the backend component of the **GoTrics** monitoring system. It collects metrics data from multiple **GoTrics nodes** and serves that data to the frontend dashboard. The server provides APIs to receive (POST) and retrieve (GET) metrics data.

---

### Features

- **Metrics Collection**: Accepts POST requests from nodes that send system metrics like CPU and memory usage.
- **Metrics Aggregation**: Stores and aggregates the received metrics data from nodes.
- **API**: Exposes a REST API for retrieving the aggregated metrics via a `GET` request.
- **CORS Support**: Allows cross-origin requests to interact with the frontend dashboard.
- **Written in Go**: Fast and lightweight backend built with Go and the Gin framework.

---

### Installation

1. Clone the repository:

```bash
git clone https://github.com/MathDesigns/gotrics-server.git
cd gotrics-server
```

2. Install dependencies
```bash
go mod tidy
```

2. Run the server
```bash
go run main.go
```
Runs on port ``8080`` by default

## License
This project is licensed under the MIT License.