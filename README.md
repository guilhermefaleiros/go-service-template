# Microservice Template with Clean Architecture in Go

This project serves as a robust microservice template using Go, Echo, Postgres, OpenTelemetry (Jaeger and Prometheus), and YAML-based configuration with Viper. It follows the principles of Clean Architecture to ensure separation of concerns and scalability.

## Features
- **Clean Architecture**: Separation of business logic, infrastructure, and delivery layers.
- **Observability**: Integrated with OpenTelemetry for tracing (Jaeger) and metrics (Prometheus).
- **Configuration Management**: YAML configuration files loaded via Viper.
- **Database Management**: Postgres integration with Flyway for migrations.
- **Middleware**: Logging and recovery implemented with Echo.
- **Dockerized Environment**: Fully containerized using Docker Compose.

## Prerequisites
- Docker and Docker Compose
- Go 1.23 or later
