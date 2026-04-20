# Changelog

## v0.10.0

### Features

* **Attack API methods** — `AttackRead`, `AttackCount`, `AttackIP`, `HitDetails`, `HitRaw` for attack inspection and hit retrieval.
* **Activity log client** — `ActivityLogRead` with object-type filtering for audit-log access.
* **Security issues API** — `GetSecurityIssues`, `GetSecurityIssueGroups`, `GetSecurityIssueGroupsCount`, and related methods for the `/v1/security_issues` endpoints; request types for filter, sort, and pagination.
* **Cursor pagination for `AttackRead`** — `AttackReadRequest.Cursor` and `Paging` fields for cursor-based iteration; response carries `Cursor` for the next page. Fields are `omitempty` — existing offset-based callers unaffected.
* **Hit block-status filter** — `HitsFilter.BlockStatus []string` (`omitempty`) for filtering hits by block outcome.

### Improvements

* **Nil-input guards** — `AttackRead`, `AttackCount`, `AttackIP`, `HitDetails`, `HitRaw`, and all new security-issues methods reject `nil` request bodies with a clear error instead of shipping `"null"` to the API.
* **IP list search encoding** — `IPListSearch` now builds nested filter query parameters correctly (`filter[rule_type][]`, `filter[list]`, `filter[query]`) matching the API's Rack-style parsing.
* **`Attack` and `ActivityLog`** added to the `API` interface so the new methods are reachable via the `api` facade.

### Breaking Changes

* **`SecurityIssueMitigations.Vpatch` type** changed to `*SecurityIssueVPatch` with `json:"vpatch,omitempty"`. The API legitimately returns issues without a vpatch mitigation — consumers that dereferenced the old value type **must add a nil check** before accessing `Vpatch.RuleID`.
* **`GetSecurityIssuesResp` and related types** use Go-idiomatic initialisms (`ID`, `URL`, `ClientID`, `HTTPMethod`, `AASMTemplate`, `RuleID`) instead of `Id`/`Url`/`ClientId`/`HttpMethod`/`AasmTemplate`/`RuleId`. JSON tags unchanged. Downstream consumers pattern-matching on Go field names must rename.

### Documentation

* Added unit tests for the new API surfaces — `attack_test.go`, `activity_log_test.go`, `get_security_issues_test.go` (~760 new lines).

## v0.9.1

### Breaking Changes

* **Removed `VulnPrefix`** from `ClientCreate` struct and `ClientInfoBody` — the field was removed from the Wallarm API. Sending it causes errors.
* **Removed `get_vulns.go`** — `Vulnerability` interface and `GetVulnRead` method removed. Unused by the provider.

### Improvements

* **`make lint` target** — added golangci-lint to GNUmakefile.
* **Test coverage** — added unit tests for all API methods: Client CRUD, Action/Hint CRUD, IP lists (deny/allow/gray), integrations (11 types), triggers, users, API specs, hits, security issues, credential stuffing, wallarm mode, overlimit settings, utils. Coverage: 25.8% → 79.5%.

## v0.9.0

### Features

* **Gzip compression** — all requests send `Accept-Encoding: gzip`, responses are decompressed transparently. ~19x reduction in response payload size.
* **Batch delete** — `HintDeleteFilter.ID` changed from `int` to `[]int`, supporting batch delete of up to 1000 rules per API call.
* **IP list cache support** — `IPListReadByRuleType` method for per-rule-type filtered reads, `IPListSearch` for targeted value lookup.
* **Hits fetch by attack_id** — fetch related hits across an attack campaign for false positive analysis.
* **Credential stuffing configs** — `CredentialStuffingConfigsRead` method for the v4 API endpoint.
* **Action API methods** — `ActionReadByHitID` for resolving hit-to-action mapping.

### Improvements

* **HTTP header handling** — request headers are now copied (not replaced), preserving Go's default transport headers.
* **APIError type** — structured error with `StatusCode` and `Body` fields, compatible with `errors.As()`.
* **Retry policy** — configurable retry for 423 (rules locked), 5xx (server error), and 429 (rate limit) with exponential backoff.
* **Pagination fix** — all paginated methods set `response.Body.Objects = nil` before each `json.Unmarshal` to prevent slice reuse bugs.

### Documentation

* **README rewrite** — updated capabilities list, added features section (retry, gzip, structured errors), updated code examples.

### Breaking Changes

* `HintDeleteFilter.ID` type changed from `int` to `[]int` — callers must wrap single IDs in `[]int{id}`.
* `ClientFields.Enabled` changed from `bool` to `*bool` (fixes `omitempty` dropping `false`).
