# Changelog

## [v0.12.0] - 2026-04-28

### Upgrade Steps

* [ACTION REQUIRED] `HintDelete` signature changed from `(*HintDelete) error` to `(*HintDelete) (*HintDeleteResp, error)`. Downstream code using `_, err := api.HintDelete(...)` keeps compiling; call sites that stored the function reference or used named returns must update.

### Breaking Changes

* **`HintDelete` returns `(*HintDeleteResp, error)`** — the new `HintDeleteResp{Status int, Body []ActionBody}` exposes the API response. Inspect `len(resp.Body) == 0` to detect the no-op path: the Wallarm API returns HTTP 200 in all cases (rule already absent, never existed, or delete blocked server-side as for counter hints), and the body is the only signal of whether anything was actually deleted.

### New Features

* **`HintUpdateV3Params` extended with the full mutable-field surface** — adds 30+ new fields (all pointer types with `omitempty`) covering common attributes (`Title`, `Active`, `Set`), per-rule mutable fields (`Mode`, `AttackType`, `Stamp`, `Regex`, `LoginRegex`, `CaseSensitive`, `CredStuffType`, `Size`, GraphQL limits, `OverlimitTime`, `Parser`, `State`, rate-limit knobs, response-header `Name`/`Values`, `FileType`), and mitigation-control nested structs (`Threshold`, `Reaction`, `EnumeratedParameters`). Existing `Comment` and `VariativityDisabled` are unchanged. Pointer + `omitempty` semantics let callers send only the fields they want to update; nil fields stay out of the wire payload.

## [v0.11.0] - 2026-04-22

> API spec client expansion with new endpoints and payload fields, plus Go-idiomatic initialism renames across the API spec and rules settings types.

### Upgrade Steps

* [ACTION REQUIRED] Rename any `wallarm.ApiSpec*` references in downstream code to `wallarm.APISpec*` (types, methods, fields, error strings) — affects `ApiSpecCreate`, `ApiSpecRead`, `ApiSpecDelete`, `ApiSpecBody`, and related identifiers across `api_spec.go` and `get_hits.go`.
* [ACTION REQUIRED] Rename `RulesSettingsResponseBody.ClientId` to `ClientID`. JSON tag `clientid` is unchanged.

### Breaking Changes

* **API spec identifiers renamed to idiomatic Go initialisms** — every `ApiSpec*` symbol is now `APISpec*` (types, methods, fields, test names, error messages). JSON tags are unchanged. Pattern-matching consumers must rename. `87b0dda`, `e2ad131`, `a803a40`
* **`RulesSettingsResponseBody.ClientId` renamed to `ClientID`** — JSON tag stays `clientid`. `4cfd2e5`

### New Features

* **New API spec endpoints** — `APISpecReadByID` (GET by id), `APISpecUpdate` (PUT), and `APISpecList` (paginated list) replace the old list-and-filter read flow. `82303f7`
* **API spec policy endpoint** — `APISpecPolicyPut` for managing API spec policies, with new `APISpecPolicy`, `APISpecPolicyCondition`, and `APISpecPolicyResp` types. `a1c7df2`
* **Extended `APISpecBody`** — added `policy`, `auth_headers`, `file`, and `format` fields, plus supporting types `APISpecAuthHeader`, `APISpecFile`, `APISpecUpdate`, and `APISpecListResp`. `a427ad9`
* **`APISpecCreate.AuthHeaders`** — auth headers for URL-based spec fetch are now settable on create, matching the update/read surface. `55e328f`

### Bug Fixes

* **`APISpecReadByID` treats HTTP 404 as `ErrNotFound`** — callers can detect deleted specs via `errors.Is(err, ErrNotFound)`. `1493674`
* **`APISpecUpdate` and `APISpecPolicyPut` treat HTTP 404 as `ErrNotFound`; `APISpecDelete` treats 404 as idempotent success** — consistent not-found semantics across the api_spec surface.
* **`APISpecUpdate` field types switched to pointers** — `Title`, `Description`, `FileRemoteURL` are now `*string` and `AuthHeaders` is `*[]APISpecAuthHeader`. `omitempty` on plain types silently dropped empty-string/empty-slice values, preventing callers from clearing fields via PUT. Pointer semantics let callers distinguish "unset" from "clear". `0c1f482`, `56d5e12`
* **`fix(api_spec_policy)`: `APISpecPolicy.Conditions` no longer has `omitempty`** — API requires the field to always be present in PUT body (empty array is valid).

### Other Changes

* **Unit tests** for `APISpecReadByID`, `APISpecUpdate`, `APISpecList`, and `APISpecPolicyPut`, plus HTTP 404 coverage for `APISpecReadByID`, `APISpecUpdate`, `APISpecPolicyPut`, and `APISpecDelete`. `6ffe6a8`
* **Test alignment** — api_spec test function names and error strings updated to match the `APISpec*` rename. `a803a40`

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
