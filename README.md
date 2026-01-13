# ⚡ GoTrics Server (ARCHIVED) V2 Soon with different architecture

THIS REPOSITORY HAS BEEN ARCHIVED

> A self-hosted, lightweight infrastructure monitoring solution designed for simplicity.

![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)
![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?logo=docker)
![License](https://img.shields.io/badge/License-MIT-green)

<div align="center">

| **[🖥️ GoTrics Server](https://github.com/MathDesigns/gotrics-server)** | **[🕵️‍♂️ GoTrics Node](https://github.com/MathDesigns/gotrics-node)** | **[📊 GoTrics Front](https://github.com/MathDesigns/gotrics-front)** |
| :---: | :---: | :---: |
| The Brain (Backend) | The Agent (Collector) | The Dashboard (UI) |

</div>

## 🧐 Why GoTrics?
Most monitoring solutions (Prometheus/Grafana) are heavy and complex to set up. **GoTrics** solves this by offering:
* **Zero-Config Agents:** Just run the binary, point it to the server, and you're done.
* **Low Footprint:** Written in Go to consume minimal resources.
* **Real-Time:** Sub-second metric updates via WebSockets.

## 🏗 Architecture
```mermaid
flowchart LR
    subgraph Infrastructure
        A[gotrics-node] --> B[CPU / RAM / Disk]
    end

    subgraph Server[GoTrics Server]
        C[(PostgreSQL)]
        D[gotrics-server]
        E[(Redis)]
        D --- C
        D --- E
    end

    subgraph Client
        F[Svelte 5 Dashboard]
    end

    B -- HTTP/WS --> D
    D -- REST --> F

    style A fill:#334155,stroke:#64748b,color:#f1f5f9
    style B fill:#475569,stroke:#94a3b8,color:#f1f5f9
    style C fill:#047857,stroke:#10b981,color:#fff
    style D fill:#4f46e5,stroke:#818cf8,color:#fff
    style E fill:#dc2626,stroke:#f87171,color:#fff
    style F fill:#f59e0b,stroke:#fbbf24,color:#1c1917
```
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
