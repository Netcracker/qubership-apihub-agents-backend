# qubership-apihub-agents-backend

A Go-based microservice backend that provides management and orchestration capabilities for Qubership APIHUB agents.
This service handles agent registration, service discovery, snapshot management, and security checks.

## API Documentation

API documentation is available in the [OpenAPI specification](docs/api/Agents-Backend-API.yaml).

## Installation

This service is designed for Kubernetes deployment and uses PostgreSQL as the database.
For deployment, use the provided [Helm chart](/helm-templates/README.md).

## Build

Just run `build_golang_binary.cmd` file.

For Docker builds, use `build_docker_image.cmd`.
