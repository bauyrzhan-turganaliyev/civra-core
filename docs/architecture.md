# Civra – Architecture Overview

## General Architecture
Civra backend follows a microservice-oriented architecture with a clear separation of concerns.

Main components:
- **Gateway** – entry point for HTTP clients
- **Economy Service** – handles resources, quotas and market logic
- **PostgreSQL** – shared persistent state

All services communicate via HTTP using JSON.

