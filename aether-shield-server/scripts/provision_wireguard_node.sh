#!/usr/bin/env bash
set -euo pipefail

: "${AETHER_WG_INTERFACE:=wg0}"
: "${AETHER_WG_ADDRESS:=10.88.0.1/24}"
: "${AETHER_WG_PORT:=51820}"
: "${AETHER_WG_PRIVATE_KEY_FILE:=/etc/wireguard/aether-server.key}"

if [[ ${EUID} -ne 0 ]]; then
  echo "run as root" >&2
  exit 1
fi

for cmd in wg ip sysctl install; do
  command -v "$cmd" >/dev/null 2>&1 || { echo "missing dependency: $cmd" >&2; exit 1; }
done

install -d -m 700 /etc/wireguard

if [[ ! -f "$AETHER_WG_PRIVATE_KEY_FILE" ]]; then
  umask 077
  wg genkey > "$AETHER_WG_PRIVATE_KEY_FILE"
  chmod 600 "$AETHER_WG_PRIVATE_KEY_FILE"
fi

PRIVATE_KEY="$(cat "$AETHER_WG_PRIVATE_KEY_FILE")"
PUBLIC_KEY="$(printf '%s' "$PRIVATE_KEY" | wg pubkey)"

CONF="/etc/wireguard/${AETHER_WG_INTERFACE}.conf"
umask 077
cat > "$CONF" <<EOF
[Interface]
Address = ${AETHER_WG_ADDRESS}
ListenPort = ${AETHER_WG_PORT}
PrivateKey = ${PRIVATE_KEY}
EOF
chmod 600 "$CONF"

cat > /etc/sysctl.d/99-aether-wireguard.conf <<EOF
net.ipv4.ip_forward=1
net.ipv6.conf.all.forwarding=1
EOF
sysctl --system >/dev/null

echo "WireGuard configuration written to $CONF"
echo "PublicKey=$PUBLIC_KEY"
echo "No peers were added. Client peer provisioning must be configured separately."
