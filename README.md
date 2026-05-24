# Flight Search & Aggregation System

A backend service that aggregates flight search results from multiple airline providers.

## Overview

This application "queries" various airline providers (Garuda Indonesia, Lion Air, Batik Air, AirAsia), normalizes their different response formats into a common domain model, and returns unified search results.

### Key Design Choices

- **Domain-driven design** — Bridging communication gaps through a shared "ubiquitous language" improving maintainability, independent modules, and building the codes/business logic around the domain.
- **Provider abstraction** — Each airline implements the `FlightProvider` interface. Adding a new provider requires only a new package and a config entry, clear separation, and isolated andimplementation details. ref: https://refactoring.guru/design-patterns/strategy
- **POST HTTP method** - While GET is the standard for searching, the POST method is used for search operations when requests are too complex or sensitive for a URL. ref: https://stackoverflow.com/a/64470896
- **Shared utilities** — Common parsing logic (time, duration, airport lookup) is extracted to `partners/utils.go` to avoid code duplication. ref: https://refactoring.guru/smells/shotgun-surgery
- **Concurrency** - Flights are fetched concurrently using goroutines + `sync.WaitGroup`. Total search latency is bounded by the slowest provider (~400ms worst case) rather than the sum of all providers (if sync).

## Project Structure
1. cmd -> entry point
2. controller -> http handlers
3. application -> business logic
4. domain -> domain model
5. partners -> flight provider integrations

## Prerequisites

- **Go 1.20** or later
- No external services, infrastrucutre, repository required (providers response are mocked with local JSON files)

## Setup & Run

1. **Clone the repository**

```bash
git clone https://github.com/febster16/flight-search-aggregator
cd flight-search-aggregator
```

2. **Install dependencies**

```bash
go mod download
```

3. **Configuration File**

config-staging.yml:
```
http:
  name: "flight-search-aggregator"
  port: 8000
```

4. **Run the server**

```bash
make run-http
```

Or manually:

```bash
ENVIRONMENT=staging go run cmd/http/main.go
```

The server starts on port **8000** (configurable in `config-staging.yml`).

## API Usage

### 1. Health Check

```bash
curl http://localhost:8000/health
```

### 2. Search Flights

```bash
curl -X POST http://localhost:8000/flights/search \
  -H "Content-Type: application/json" \
  -d '{
    "origin": "CGK",
    "destination": "DPS",
    "departure_date": "2025-12-15",
    "return_date": null,
    "passengers": 1,
    "cabin_class": "economy"
  }'
```

## Future Improvements / TODOs

1. **FlightProvider injection using `map[string]partners.FlightProvider`** for named provider access, enabling per-provider retry policies and circuit breakers.
2. **Search request validation** — enforce required fields and value constraints (e.g., `github.com/go-playground/validator`).
3. **Baggage parsing** — AirAsia's free-text baggage notes need a more nuanced parser to extract structured carry-on/checked values.
4. **Airport city lookup** — replace internal dictionary with a full library (e.g., `mmcloughlin/openflights`) for production coverage.
5. **Go worker pool** — implement worker pool to better handle goroutines under high traffic.
6. **Config-driven provider toggling** — load enabled providers from config YAML for enable/disable without code changes.
7. **Caching layer** — cache or `single flight` provider responses based on filter to reduce redundant calls.
8. **"Best value" scoring** — rank results by weighted combination of price and convenience.
9. **Timeout handling** — add context deadline for provider calls to prevent unbounded waits.
