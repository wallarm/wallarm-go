# wallarm-go

[![build](https://github.com/416e64726579/wallarm-go/workflows/Go/badge.svg)](https://github.com/416e64726579/wallarm-go/actions?query=workflow%3AGo)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/416e64726579/wallarm-go)](https://pkg.go.dev/github.com/416e64726579/wallarm-go)
[![codecov](https://codecov.io/gh/416e64726579/wallarm-go/branch/master/graph/badge.svg)](https://codecov.io/gh/416e64726579/wallarm-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/416e64726579/wallarm-go?style=flat-square)](https://goreportcard.com/report/github.com/416e64726579/wallarm-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://github.com/416e64726579/wallarm-go/blob/master/LICENSE)

## Table of Contents
- [Install](#install)
- [Getting Started](#getting-started)
- [License](#license)

> **Note**: This library is in active development and highly suggested to use carefully.

A Go library for interacting with
[Wallarm API](https://apiconsole.eu1.wallarm.com). This library allows you to:

* Manage applications
* Manage nodes
* Manage integrations
* Manage triggers
* Manage users
* Manage the blacklist
* Switch the WAF/Scanner/Active Threat Verification modes
* Inquire found vulnerabilities

## Install

You need a working Go environment

```sh
go get github.com/416e64726579/wallarm-go
```

## Getting Started

The sample code could be similar 

```go
package main

import (
	"log"
	"net/http"
	"os"

	wallarm "github.com/416e64726579/wallarm-go"
)

func main() {

	wapiHost, exist := os.LookupEnv("WALLARM_API_HOST")
	if !exist {
		wapiHost = "https://api.wallarm.com"
	}
	wapiUUID, exist := os.LookupEnv("WALLARM_API_UUID")
	if !exist {
		log.Fatal("ENV variable WALLARM_API_UUID is not present")
	}
	wapiSecret, exist := os.LookupEnv("WALLARM_API_SECRET")
	if !exist {
		log.Fatal("ENV variable WALLARM_API_SECRET is not present")
	}

	authHeaders := make(http.Header)
	authHeaders.Add("X-WallarmAPI-UUID", wapiUUID)
	authHeaders.Add("X-WallarmAPI-Secret", wapiSecret)

	// Construct a new API object
	api, err := wallarm.New(wallarm.UsingBaseURL(wapiHost), wallarm.Headers(authHeaders))
	if err != nil {
		log.Print(err)
	}

	// Fetch user details
	u, err := api.UserDetails()
	if err != nil {
		log.Print(err)
	}
	// Print user specific data
	log.Println(u.Body)

	// Change global WAF mode to monitoring
	mode := wallarm.ClientUpdate{
		Filter: &wallarm.ClientFilter{
			ID: 1,
		},
		Fields: &wallarm.ClientFields{
			Mode: "monitoring",
		},
	}
	c, err := api.ClientUpdate(&mode)
	if err != nil {
		log.Print(err)
	}
	// Print client data
	log.Println(c)

	// Create a trigger when the number of attacks more than 1000 in 10 minutes
	filter := wallarm.TriggerFilters{
		ID:       "ip_address",
		Operator: "eq",
		Values:   []interface{}{"2.2.2.2"},
	}

	var filters []wallarm.TriggerFilters
	filters = append(filters, filter)

	action := wallarm.TriggerActions{
		ID: "send_notification",
		Params: wallarm.TriggerActionParams{
			IntegrationIds: []int{5},
		},
	}

	var actions []wallarm.TriggerActions
	actions = append(actions, action)

	triggerBody := wallarm.TriggerCreate{
		Trigger: &wallarm.TriggerParam{
			Name:       "New Terraform Trigger Telegram",
			Comment:    "This is a description set by Terraform",
			TemplateID: "attacks_exceeded",
			Enabled:    true,
			Filters:    &filters,
			Actions:    &actions,
		},
	}

	triggerResp, err := api.TriggerCreate(&triggerBody, 1)
	if err != nil {
		log.Print(err)
	}
	// Print trigger metadata
	log.Println(triggerResp)
}
```

# License

[MIT](LICENSE) licensed
