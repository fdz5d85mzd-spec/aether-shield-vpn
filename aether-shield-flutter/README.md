# Mobile integration status

This directory contains native integration scaffolds only.

- Android: VpnService can establish a TUN interface, but no Rust packet loop is connected yet.
- iOS: PacketTunnelProvider can apply tunnel settings, but packetFlow processing and Rust FFI are not implemented.
- iOS NetworkExtension entitlement, signing, App Group/configuration, and on-device validation remain required.
- Neither scaffold is evidence of a working VPN connection.
