use crate::config::TunnelConfig;

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum EngineState {
    Unsupported,
    Ready,
    Active,
    Failed,
}

#[derive(Debug, Clone, Copy, Default, PartialEq, Eq)]
pub struct EngineTelemetry {
    pub packets_from_tun: u64,
    pub packets_to_tun: u64,
    pub encrypted_bytes_sent: u64,
    pub encrypted_bytes_received: u64,
}

pub trait PacketEngine: Send {
    fn configure(&mut self, config: &TunnelConfig) -> Result<(), &'static str>;
    fn start(&mut self, tun_fd: i32) -> Result<(), &'static str>;
    fn stop(&mut self);
    fn state(&self) -> EngineState;
    fn telemetry(&self) -> EngineTelemetry;
}

#[derive(Debug, Default)]
pub struct UnsupportedEngine;

impl PacketEngine for UnsupportedEngine {
    fn configure(&mut self, _config: &TunnelConfig) -> Result<(), &'static str> {
        Err("WireGuard packet engine unavailable")
    }

    fn start(&mut self, _tun_fd: i32) -> Result<(), &'static str> {
        Err("WireGuard packet engine unavailable")
    }

    fn stop(&mut self) {}

    fn state(&self) -> EngineState {
        EngineState::Unsupported
    }

    fn telemetry(&self) -> EngineTelemetry {
        EngineTelemetry::default()
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::config::{PeerConfig, TunnelConfig};

    #[test]
    fn unsupported_engine_never_reports_traffic_or_active() {
        let mut engine = UnsupportedEngine;
        let config = TunnelConfig {
            private_key: "runtime-only".into(),
            addresses: vec!["10.88.0.10/32".into()],
            dns_servers: vec![],
            mtu: 1420,
            peer: PeerConfig {
                public_key: "server-public".into(),
                endpoint: "example.invalid:51820".into(),
                allowed_ips: vec!["0.0.0.0/0".into()],
                persistent_keepalive_seconds: Some(25),
            },
        };
        assert!(engine.configure(&config).is_err());
        assert!(engine.start(7).is_err());
        assert_eq!(engine.state(), EngineState::Unsupported);
        assert_eq!(engine.telemetry(), EngineTelemetry::default());
    }
}
