use crate::config::TunnelConfig;
use crate::transport::{PacketTransport, TransportState};

#[derive(Debug)]
pub struct WireGuardTransport {
    state: TransportState,
    config: TunnelConfig,
}

impl WireGuardTransport {
    pub fn new(config: TunnelConfig) -> Result<Self, &'static str> {
        config.validate()?;
        Ok(Self {
            state: TransportState::Ready,
            config,
        })
    }

    pub fn endpoint(&self) -> &str {
        &self.config.peer.endpoint
    }
}

impl PacketTransport for WireGuardTransport {
    fn state(&self) -> TransportState {
        self.state
    }

    fn start(&mut self) -> Result<(), &'static str> {
        self.state = TransportState::Failed;
        Err("WireGuard packet engine not integrated")
    }

    fn stop(&mut self) {
        self.state = TransportState::Ready;
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::config::PeerConfig;

    #[test]
    fn adapter_never_fakes_activation() {
        let config = TunnelConfig {
            private_key: "runtime-only".into(),
            addresses: vec!["10.88.0.2/32".into()],
            dns_servers: vec![],
            mtu: 1420,
            peer: PeerConfig {
                public_key: "peer-runtime-only".into(),
                endpoint: "vpn.example.invalid:51820".into(),
                allowed_ips: vec!["0.0.0.0/0".into()],
                persistent_keepalive_seconds: None,
            },
        };
        let mut transport = WireGuardTransport::new(config).unwrap();
        assert_eq!(transport.state(), TransportState::Ready);
        assert!(transport.start().is_err());
        assert_eq!(transport.state(), TransportState::Failed);
    }
}
