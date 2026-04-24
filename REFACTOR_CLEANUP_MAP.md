# Refactor Cleanup Map

This file is the source-of-truth for the aggressive simplification pass.

## Runtime Primary Path (kept)

- `cmd/server/main.go`
- `internal/handlers/*`
- `internal/domain/services/*` (non-DDD services)
- `internal/domain/repositories/*` (interfaces used by active handlers/services)
- `public/html/parking.html` + `public/js/*` (parking UI)
- `scripts/test_parking_apis.sh` and `scripts/test_avp_apis.sh`

## Decoupled From Default Runtime

- `internal/integration/ddd_service.go` and related DDD integration stack are no longer started from `main`.
- Result: one runtime entry path by default, no parallel architecture bootstrapping at process start.

## Merge Targets

- Parking distance calculation and mock object shaping are centralized under parking domain services.
- AVP task flow remains in a single dispatch service; handler now only orchestrates request/response mapping.
- API test scripts share one common helper library.

## Delete Targets (safe cleanup candidates, next pass)

- DDD demo-specific repositories/factories not referenced by runtime routes.
- Duplicate documentation sections that repeat endpoint lists and test steps verbatim.
- Any helper methods with zero call sites after service consolidation.
