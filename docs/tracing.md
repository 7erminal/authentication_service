# Request correlation and safe logging

## Implemented

- Service name: `authentication_service`. Go module version is 1.24; logging uses the standard-library `log/slog`. The Docker builder now uses Go 1.24 as well. No dependencies were added or upgraded.
- `internal/trace` is a **local compatibility boundary**, not a copy claimed to come from the gateway. No shared trace package was found in this repo, dependencies, or the inspected gateway branches (`main`, `ecommerce`, `shopmallafrik`). Replace this boundary when a shared module is published.
- Accept `X-Session-Id`, `X-Trace-Id`, `X-Request-Id` only as a single header value with 8–64 ASCII characters from `[A-Za-z0-9_-]`. Reject the entire invalid value, including whitespace, commas, Unicode, control characters, and ambiguous repeated headers. Invalid values are never logged or echoed.
- Generate a cryptographically random request ID if missing/invalid. A missing/invalid trace ID falls back to that request ID. A missing/invalid session ID remains empty. Missing headers never reject an otherwise valid request.
- A `BeforeRouter` filter stores/accesses the IDs through the standard request context and writes request/trace response headers. The outer HTTP middleware initializes the same context earlier so static responses, oversized requests, unknown methods, and other pre-router exits are covered too. Initialization is idempotent and never replaces an established request ID.
- The filter is installed before CORS. CORS allows all three tracing headers and exposes the request/trace response headers.
- JSON logs go to **stdout** instead of the old application file. Deployment log collection/rotation must collect stdout. Every application event includes `service`, `session_id`, `trace_id`, and `request_id`; startup/non-request events use empty ID fields.
- Application severity is `LOG_LEVEL`, falling back to `loglevel` in `conf/app.conf`, then `INFO`. Supports standard slog levels (`DEBUG`, `INFO`, `WARN`, `ERROR`). Invalid levels fail startup without printing the supplied value. Access events remain enabled at INFO even when application logging is configured at WARN/ERROR.
- One deferred `http access` event includes method, redacted path, route, status, `duration_ms`, and response-body bytes. The matched route template is the redacted path; static/unmatched/pre-routing responses use `<unmatched>`. No query strings or raw user-supplied paths are logged.
- Beego built-in access/static logs are explicitly disabled. Unscoped Beego diagnostics, standard-library free-form logs, and ORM debug output are suppressed: those APIs cannot safely attach request context and can expose URLs, SQL arguments, or credentials. Safe application error, outbound-failure, access, panic, and startup events replace them. Do not re-enable HTTP/ORM debug dumps.
- Panics emit one correlated error event containing stack function/file/line locations, **not** panic values or argument dumps. Clients receive HTTP 500 and `{"error":{"code":"internal_error","message":"Internal server error","request_id":"..."}}`. Beego's deliberate `ErrAbort` sentinel is not treated as a panic.
- Ordinary responses buffer up to 1 MiB of body data, allowing even a panic after a partial JSON write to replace the response with a clean 500. Partial cookies, redirects, and representation headers are discarded; CORS/correlation headers survive. Larger responses, `Flush`, and `Hijack` commit the network response. HTTP cannot change status after that point: a later panic logs `response_committed=true`, does not append error JSON to the partial body, and still emits one access event. There are no application streaming handlers identified in this service.

IDs are metadata only: they are never used for auth decisions, primary keys, cache keys, deduplication, or idempotency.

## Error compatibility

`tracedController.respondError(status, code, message)` delegates to the common JSON error writer and always includes the request ID. Beego framework errors use controller-style error handlers (not callback handlers, which commit headers too early). Their existing HTTP statuses remain unchanged; their old HTML/plain responses become JSON with codes such as `http_404`.

Existing application failures are **not** replaced wholesale with the new envelope:

- Keep HTTP status, `StatusCode`, `StatusDesc`, `Result`, `Value`, and other existing values.
- Add only `request_id` to object-shaped failures (`StatusCode >= 300`, HTTP status >= 400, or `IsValid == false`). This includes legacy business failures returned as HTTP 200.
- Keep success response schemas unchanged.
- Keep string-shaped failures as strings. They have request correlation in response headers, but not in the body. Converting these to objects needs separate compatibility approval.
- Existing OAuth transport failures previously caused nil-response panics/500s; they now return the same 500 through the safe error helper instead.

The inspected gateway unmarshals authentication responses into typed DTOs and branches on `StatusCode`/reads `StatusDesc` and `Result`. Those fields are preserved; Go's normal JSON decoder ignores the additive `request_id`. Full gateway integration testing is still required before rollout. A strict non-Go client may need its response schema to allow this additional field.

## Propagation coverage

### Request flows migrated

All currently active request->model call sites use context-aware variants, including:

- User/customer login and credential verification.
- Access/refresh-token creation, lookup, verification, refresh, expiry, and revocation.
- Credential create/update and password-change/reset flows.
- OTP, activation-code, invitation-token, and password-reset-link storage.

All eight record model files use ORM `InsertWithCtx`, `ReadWithCtx`, `OneWithCtx`, `AllWithCtx`, `UpdateWithCtx`, and `DeleteWithCtx` as appropriate. They do not use deprecated `QueryTableWithCtx`, which would not propagate context to execution. Successful inserts log the internal numeric `record_id` with `record_type` and tracing IDs; tokens and entire model values are never logged. Compatibility model functions without a context remain as wrappers using `context.Background()`; no active controller calls these wrappers.

All internal HTTP helpers (user/customer lookup by ID or username, role lookup, user-password update) use `SendRequestContext` with `Internal: true`. They forward the **same** three IDs, omitting absent session IDs. Caller-supplied trace headers are replaced with context values, never used to generate new IDs. Redirects retain IDs only on the original origin; cross-origin redirects strip them.

The bearer-token user lookup no longer constructs an unassociated empty controller; it calls `GetUserContext` with the originating context.

Google OAuth requests receive context/cancellation, but **not tracing headers**. External requests and redirects have correlation headers removed. Request serialization, auth headers, and normal application payloads are otherwise preserved. No live Google/customer-service calls were made during validation.

### Background work and remaining limitations

- No cache client, queue producer, or queue consumer was found. No artificial queue abstraction was added.
- The existing password-reset email goroutine receives a detached context containing only sanitized IDs; it does not retain the pooled Beego controller or arbitrary request values. Email completion/failure logs use that context.
- `trace.IDs` can be serialized into future message metadata, restored through `trace.WithIDs`, and used without generating a new request ID. Metadata serialization/restoration is tested.
- SMTP still uses `net/smtp.SendMail`, which has no context cancellation/deadline API. It carries IDs for logs only, not HTTP headers or email content. A deadline-aware SMTP implementation/reliable job queue is separate work.
- Non-request callers of compatibility wrappers (`SendRequest`, model functions without `Context`, `GetUserFromBearerToken`, `Register`, `SendEmail`, `SendEmailNew`) have empty tracing fields unless migrated to the context-aware entry points.
- DB registration/Ping/timezone discovery happen at startup, without request context.

## Remaining non-migrated log calls

| Location | Treatment |
| --- | --- |
| `loader/main.go` schema-tool error `fmt.Fprintf(os.Stderr, ...)` | Unchanged standalone Atlas/schema-loader diagnostic; not part of the serving process. Its stdout must remain SQL, not JSON logs. |
| Beego/global standard logger/ORM internals | Unscoped free-form output intentionally suppressed after logging configuration. Framework/library initialization before `main` is outside the request lifecycle. |
| Active controllers, token/email/HTTP helpers, and models | No remaining legacy `logs.*`/`fmt.Print*` calls. Unsafe calls were removed, not wrapped with their original payloads. |

The existing tracked executable is unchanged. Deploy a newly built binary, not that historical executable.

## Database proposal — NOT applied

Candidate authentication-owned tables:

- `access_tokens`
- `refresh_tokens`
- `customer_access_tokens`
- `customer_refresh_tokens`
- `customer_credentials`
- `user_tokens`
- `user_otps`
- `activation_codes`

Proposed change per table, pending **separate schema approval**:

```sql
-- Proposal only. Repeat for each approved table; do not run automatically.
ALTER TABLE access_tokens
  ADD COLUMN trace_id VARCHAR(64) NULL,
  ADD COLUMN request_id VARCHAR(64) NULL,
  ADD COLUMN session_id VARCHAR(64) NULL;
```

These columns would be observability-only, populated from the originating creation request. No indexes, constraints, primary-key changes, auth/key use, or idempotency behavior are proposed. No migration file, model field, or live database change was made. Broader historical customer/user/reference/DTO-derived tables are excluded unless ownership is explicitly confirmed. Plan rollout/backfill/retention policy before adding persistence.

## Example logs

Illustrative JSON lines (not real customer data):

```json
{"time":"2026-11-27T16:00:00Z","level":"INFO","msg":"record created","service":"authentication_service","session_id":"session_123","trace_id":"trace_123","request_id":"request_123","record_type":"access_tokens","record_id":42}
{"time":"2026-11-27T16:00:00Z","level":"INFO","msg":"http access","service":"authentication_service","session_id":"session_123","trace_id":"trace_123","request_id":"request_123","method":"POST","path":"/v1/auth/login/token","route":"/v1/auth/login/token","status":200,"duration_ms":12.4,"bytes":512}
{"time":"2026-11-27T16:00:00Z","level":"WARN","msg":"application error","service":"authentication_service","session_id":"","trace_id":"trace_123","request_id":"request_123","code":605,"http_status":200}
```

A request with no gateway headers receives generated request/trace IDs and an empty session field; it is not rejected.

## Validation

Run with Go 1.24+:

```sh
go test ./...
go test -race ./...
go vet ./...
go build -o /tmp/authentication-service .
```

Tests cover valid/boundary/invalid/repeated/missing/injected headers; generated IDs; independent fallbacks; metadata restoration; success/error/panic/abort/early/pre-router responses; partial-write recovery; framework JSON content type; single access events (including ERROR log level); buffer boundaries; concurrent-request isolation; legacy application schema/status compatibility; real Beego endpoint dispatch; all internal HTTP helper forwarding; redirect trust boundaries; external-provider exclusion; HTTP/SQL cancellation; SQL driver context identity; business-ID logging; and sensitive values absent from logs.

Validation in this sandbox used Go 1.24.3 and a temporary GitHub-backed module cache because the public Go module proxy was unreachable. Existing `go.sum` checks were retained. `go.mod` and `go.sum` are unchanged. Test SQL/HTTP servers are local fixtures: no live DB, SMTP, external provider, or gateway deployment was used. Build/cache artifacts were kept outside Git.

## Google OAuth deployment configuration

GitHub push protection flagged pre-existing hard-coded Google OAuth credentials
in the touched OAuth controller. Those literals have been removed from this PR.
Configure these environment variables through your deployment secret manager
before deploying:

- `GOOGLE_OAUTH_CLIENT_ID`: used by the Google authorization URL and token exchange.
- `GOOGLE_OAUTH_CLIENT_SECRET`: used only by the token exchange.

There is no hard-coded fallback. Missing configuration returns the safe correlated
500 error on the affected Google OAuth endpoint; other authentication endpoints
remain available. Values are never logged. No credential values were added to
configuration files or documentation. Rotate the previously exposed Google OAuth
credentials; removing them from the new commit does not remove historical copies.

## Open follow-ups

1. Supply the shared gateway trace module path/version when available.
2. Approve a future conversion of legacy string-shaped errors if body-level request IDs are required there too; coordinate strict client schemas before rollout.
3. Separately approve exact tables, retention, and rollout for trace-column persistence. No schema work is assumed approved by the code implementation.
4. Confirm stdout collection/retention in deployment; the former application-file logger is no longer written.
5. Rotate the previously exposed Google OAuth credentials and configure the new environment variables before deployment. DB/SMTP/signing/encryption credentials still need rotation and externalization in a separate security change; those values were not modified here.
6. Consider deadline-aware/reliable email delivery. Post-commit streaming panic status cannot be rewritten by HTTP; no such application handlers currently exist.
