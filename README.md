# AETHER SHIELD VPN

AETHER SHIELD is being built as a testable VPN platform foundation with a React control studio, Go control-plane API, Rust native core, and native Android/iOS tunnel adapters.

> **Status:** foundation/scaffold. This repository does **not** yet claim a verified end-to-end VPN tunnel, deployed RAM-node network, DPI bypass, post-quantum protection, or completed leak protection.

## Implementation status

| Area | Status |
|---|---|
| Web Control Studio | PARTIAL — supplied UI is being integrated |
| Go control-plane | SCAFFOLD |
| Rust tunnel core | SCAFFOLD |
| Protocol policy engine | SCAFFOLD |
| Android VpnService | SCAFFOLD |
| iOS PacketTunnelProvider | SCAFFOLD |
| Real node infrastructure | REQUIRES INFRASTRUCTURE |
| End-to-end VPN tunnel | NOT IMPLEMENTED |
| Post-quantum key exchange | NOT IMPLEMENTED |
| Security/leak audit | NOT IMPLEMENTED |

Production UI must never silently substitute simulated telemetry or PASS results when real services are unavailable.

## Repository layout

- `aether-shield-web/` — React + TypeScript Control Studio.
- `aether-shield-server/` — Go control-plane/API.
- `aether-shield-core/` — Rust native core and deterministic protocol policy.
- `aether-shield-flutter/` — mobile/native integration scaffold.
- `.github/workflows/ci.yml` — build and test gates.

## Security principles

No production private keys, access tokens, certificates, or credentials belong in git. Cryptographic claims must be backed by implemented and tested primitives. Node availability and audit PASS states must come from real evidence, not timers or random values.
