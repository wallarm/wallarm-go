# wallarm-go

[![build](https://github.com/wallarm/wallarm-go/workflows/Go/badge.svg)](https://github.com/wallarm/wallarm-go/actions?query=workflow%3AGo)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/wallarm/wallarm-go)](https://pkg.go.dev/github.com/wallarm/wallarm-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://github.com/wallarm/wallarm-go/blob/master/LICENSE)

A Go client library for the [Wallarm API](https://docs.wallarm.com/api/overview/). Used by the [Terraform provider for Wallarm](https://github.com/wallarm/terraform-provider-wallarm).

## Capabilities

* Rules — create, read, update, delete rules (hints) with paginated bulk fetch
* IP lists — manage allowlist, denylist, and graylist entries (subnets, countries, datacenters, proxy types)
* Integrations — email, Slack, Splunk, PagerDuty, OpsGenie, Datadog, and more
* Triggers — threshold-based alerting and reaction rules
* Applications — create and manage application pools
* Nodes — filtering node lifecycle management
* Tenants — multi-tenant account management
* Users — user CRUD and role management
* Actions — rule scope management (action conditions)
* Hits — fetch detected hits for false positive analysis
* Attacks — fetch detected attacks, hit details, and raw hit payloads
* Activity log — audit-log access with object-type filtering
* Security issues — detected vulnerabilities and attack-surface issues grouped by type and severity
* API specs — API specification management
* Credential stuffing — credential stuffing detection configs
* Global mode — filtration mode (monitoring/blocking) management

## Features

* **Automatic retry** — transient errors are retried automatically:

  | Status | Delay | Max retries |
  |--------|-------|-------------|
  | 423 (Rules locked) | 5s fixed | 12 |
  | 5xx (Server error) | 10s fixed | 12 |
  | 429 (Rate limit) | Exponential backoff | 12 |

* **Gzip compression** — requests include `Accept-Encoding: gzip` for reduced response sizes
* **Structured errors** — `APIError` type with `StatusCode` and `Body` fields, compatible with `errors.As()`
* **Configurable** — custom HTTP client, base URL, retry policy, user agent, and headers via functional options

## Install

```sh
go get github.com/wallarm/wallarm-go
```

## Getting Started

```go
package main

import (
	"log"
	"net/http"
	"os"

	wallarm "github.com/wallarm/wallarm-go"
)

func main() {
	host := os.Getenv("WALLARM_API_HOST")
	if host == "" {
		host = "https://api.wallarm.com"
	}
	token := os.Getenv("WALLARM_API_TOKEN")
	if token == "" {
		log.Fatal("WALLARM_API_TOKEN is required")
	}

	headers := http.Header{}
	headers.Set("X-WallarmAPI-Token", token)

	// Create client with token auth, retry policy, and custom base URL.
	api, err := wallarm.New(
		wallarm.UsingBaseURL(host),
		wallarm.Headers(headers),
		wallarm.UsingRetryPolicy(12, 1, 30),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Fetch user details
	user, err := api.UserDetails()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Authenticated as client %d", user.Body.ClientID)

	// Read all rules for a client
	resp, err := api.HintRead(&wallarm.HintRead{
		Limit:   500,
		OrderBy: "updated_at",
		Filter:  &wallarm.HintFilter{Clientid: []int{user.Body.ClientID}},
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Found %d rules", len(*resp.Body))
}
```

## Error Handling

API errors are returned as `*wallarm.APIError` with status code and response body:

```go
import "errors"

var apiErr *wallarm.APIError
if errors.As(err, &apiErr) {
	log.Printf("API error %d: %s", apiErr.StatusCode, apiErr.Body)
}
```

## License

[MIT](LICENSE)
