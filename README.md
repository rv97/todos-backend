# todos-backend

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/yourusername/todos-backend/blob/main/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/yourusername/todos-backend)](https://goreportcard.com/report/github.com/yourusername/todos-backend)

Welcome to **todos-backend**, the backend service for your Todos application. This project is built using Go Lang, PostgreSQL, Kafka, Docker, and Kowl.

## Features

- **RESTful API**: Provides a RESTful API for managing todos.
- **Database**: Uses PostgreSQL as the database to store todos.
- **Message Queuing**: Utilizes Kafka for asynchronous processing and real-time updates.
- **Dockerized**: Easily deploy and run the application in a Docker container.
- **Monitoring**: Monitor your Kafka cluster using Kowl.
- **Load Testing**: Conduct load testing with k6 for performance evaluation.

## Getting Started

Follow these steps to get the **todos-backend** up and running:

### Prerequisites

- Go Lang: [Installation Guide](https://golang.org/doc/install)
- PostgreSQL: [Installation Guide](https://www.postgresql.org/download/)
- Docker: [Installation Guide](https://docs.docker.com/get-docker/)
- Kafka: [Installation Guide](https://kafka.apache.org/quickstart)
- Kowl: [Installation Guide](https://github.com/cloudhut/kowl)
- k6: [Installation Guide](https://k6.io/docs/getting-started/installation/)
- **OpenTelemetry**: OpenTelemetry is used for tracing and monitoring. You can find installation instructions and documentation [here](https://opentelemetry.io/docs/go/getting-started/).

### Installation

1. Clone the repository:

```bash
git clone https://github.com/rv97/todos-backend.git
cd todos-backend