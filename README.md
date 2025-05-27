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


## Tool Architecture

![image](https://github.com/user-attachments/assets/c3fc6258-163b-4a74-b5a8-bd71cde5eabc)




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
