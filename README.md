# GoBank

A modern banking platform that leverages Go's performance for microservices, Docker for containerization, Kubernetes for orchestration, and Kafka for real-time event streaming. Built with scalability, resilience, and cloud-native principles in mind.

![image](https://github.com/user-attachments/assets/c16d8536-2252-4b95-a221-d2fa93866264)



## Overview

This project demonstrates a scalable and resilient banking platform using:
- **Go**: High-performance backend services
- **Docker**: Containerization for consistent deployment
- **Kubernetes**: Orchestration and scaling
- **Kafka**: Message streaming for asynchronous communication

## Features

- **Card Service**: Manage card operations
- **Payment Service**: Handle payment processing
- **Customer Service**: Manage customer data
- **Containerization**: Docker ensures consistent environments
- **Orchestration**: Kubernetes manages deployment, scaling, and resilience
- **Message Streaming**: Kafka enables asynchronous communication between services

## Prerequisites

- Docker
- Kubernetes cluster (e.g., Minikube, Kind, or cloud provider (I used Docker Desktop Kubernetes))
- Kafka cluster
- Go 1.16+


## Architecture Schema

flowchart LR
    subgraph CLIENT [Client]
        A[Client]
    end

    subgraph Ingress_Layer [Ingress Layer]
        B[NGINX Ingress Controller]
        C[KrakenD API Gateway]
    end

    subgraph Services [Microservices]
        D1[Customer Service<br><sub><code>gRPC</code> via Dockerfile</sub>]
        D2[Card Service<br><sub><code>gRPC</code> via Dockerfile</sub>]
        D3[Payment Service<br><sub><code>gRPC</code> via Dockerfile</sub>]
    end

    subgraph Messaging [Kafka Ecosystem]
        E[Kafka Message Queue]
        F[Zookeeper<br><sub>Coordination, Config, Cluster Mgmt</sub>]
    end

    subgraph Deployments [Kubernetes Deployments]
        G1[Service Deployments]
        G2[KrakenD Deployment]
        G3[Init DB Deployment]
        G4[PostgreSQL Deployment]
        G5[Kafka & Zookeeper Deployment]
    end

    subgraph Docker [Docker Compose Services]
        H1[PostgreSQL]
        H2[CustomerService]
        H3[CardService]
        H4[PaymentService]
        H5[KrakenD API Gateway]
        H6[Kafka Message Queue]
    end

    %% Connections
    A --> B --> C
    C --> D1
    C --> D2
    C --> D3

    D1 -->|Decrease from Personal Account| E
    D2 -->|Decrease from Card| E
    D3 -->|Payment Success Event| E

    F -->|Cluster Management| E

    style A fill:#ffffff,stroke:#000000,color:#000000
    style B fill:#00cc99,color:#ffffff
    style C fill:#00aaff,color:#ffffff
    style D1 fill:#6666ff,color:#ffffff
    style D2 fill:#6666ff,color:#ffffff
    style D3 fill:#6666ff,color:#ffffff
    style E fill:#cc00cc,color:#ffffff
    style F fill:#ff66cc,color:#ffffff
    style G1 fill:#6699cc,color:#ffffff
    style G2 fill:#6699cc,color:#ffffff
    style G3 fill:#6699cc,color:#ffffff
    style G4 fill:#6699cc,color:#ffffff
    style G5 fill:#6699cc,color:#ffffff
    style H1 fill:#3399ff,color:#ffffff
    style H2 fill:#3399ff,color:#ffffff
    style H3 fill:#3399ff,color:#ffffff
    style H4 fill:#3399ff,color:#ffffff
    style H5 fill:#3399ff,color:#ffffff
    style H6 fill:#3399ff,color:#ffffff




## Quick Start

1. **Clone the repository:**
   ```bash
   git clone https://github.com/ozturkeniss/gobank.git
   cd gobank
   ```

2. **Build and run with Docker:**
   ```bash
   docker-compose up -d
   ```

3. **Deploy to Kubernetes:**
   ```bash
   kubectl apply -f deployments/
   ```

4. **Access the API:**
   ```bash
   curl http://localhost:8085/api/v1/cards/1
   ```

## Project Structure
gobank/
├── api/proto
│ └── protocol buffers
├── cmd/
│ └── server/
├── deployments/
│ ├── krakend-config.yaml
│ └── krakend-deployment.yaml
├── internal/
│ ├── db/
│ ├── service/
│ └── payment/
├── docker/
├── .gitignore
└── README.md


## Configuration

- **KrakenD Config**: `deployments/krakend-config.yaml`
- **Kubernetes Deployment**: `deployments/krakend-deployment.yaml`
- **Environment Variables**: `.env` (not tracked in git)

## Development

- **Local Development:**
  ```bash
  go run cmd/server/main.go
  ```

- **Testing:**
  ```bash
  go test ./...
  ```

## License

MIT
