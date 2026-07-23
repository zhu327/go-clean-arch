# Structural Rules

Non-negotiables for this review. Prefer deleting complexity over rearranging it.

## Non-negotiable

1. **Do not let a PR push a file from under 1k lines to over 1k** without a strong structural reason. Prefer extract/split first.
2. **No spaghetti growth** — new ad-hoc conditionals, scattered special cases, or one-off branches in unrelated paths are design problems, not nits.
3. **Bias to cleaner design** — same behavior, meaningfully simpler structure wins. Moving complexity is not enough.
4. **Direct over magical** — flag brittle/ad-hoc "magic", thin identity wrappers, pass-through helpers that buy no clarity.
5. **Type/boundary cleanliness** — question unnecessary `any` / untyped escape hatches / cast-heavy contracts when a typed model would simplify control flow.
6. **Canonical home** — feature logic belongs in the owning layer/package; reuse existing helpers over near-duplicates.
7. **Orchestration** — avoid pointless serialization of independent work; avoid partial-update flows that leave state half-applied.

## Preferred remedies

| Instead of... | Prefer... |
|---------------|-----------|
| Polishing indirection | Deleting the layer |
| Centralizing conditionals | Reframing state so conditionals disappear |
| Helper wrapping mess | Changing ownership so the feature fits naturally |
| Special-case branches | Simpler default with fewer exceptions |
| Moving complexity | Restructuring so it no longer exists |

## Removal candidates

- **Safe delete now**: no references (including dynamic), no external consumers, tests/docs updated.
- **Defer**: active consumers, needs migration/telemetry/sign-off — list preconditions, steps, verification, rollback.
