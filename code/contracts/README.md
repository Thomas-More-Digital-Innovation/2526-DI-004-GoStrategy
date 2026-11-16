# Contracts

Shared JSON fixtures for rules the frontend mirrors from the backend.
Both the Go tests (`*contract_test.go`) and the Vitest spec (`code/frontend/src/lib/contracts.spec.ts`) run their own implementation against these files, so drift on either side fails CI.
Change a fixture only together with both implementations.
