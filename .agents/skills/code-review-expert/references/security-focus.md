# Security Focus

Assume general OWASP knowledge. Only escalate what shows up in the diff.

## AuthZ / multi-tenant

- New or changed endpoints missing auth/RBAC guards.
- Missing tenant/ownership checks on read/write (IDOR).
- Trusting client-provided roles, tenant IDs, or resource IDs.

## Data access

- Query/command construction via string concat instead of parameterized APIs.
- User-controlled paths in file I/O (`../`) or URLs reaching internal services (SSRF).
- User input into shell / process execution.

## Secrets & errors

- Credentials, tokens, or PII in code, logs, or API error payloads.
- Do not leak internals (stack traces, SQL, host details) to clients — use the project's boundary error type.

## Concurrency & resources

- Check-then-act on shared state or DB rows without lock/transaction/version.
- Unbounded fan-out, missing timeouts on external calls, large in-memory buffers on request path.
- Background-task / async errors discarded or never joined.

## When reporting

State **exploitability** and **impact**. Skip speculative supply-chain/CVE audits unless the diff changes deps or trust boundaries.
