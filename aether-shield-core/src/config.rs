use serde::{Deserialize, Serialize};

#[derive(Debug, Clone, Serialize, Deserialize, PartialEq, Eq)]
pub struct PeerConfig {
    pub public_key: String,
    pub endpoint: String,
    pub allowed_ips: Vec<String>,
    pub persistent_keepalive_seconds: Option<u16>,
}

#[derive(Debug, Clone, Serialize, Deserialize, PartialEq, Eq)]
pub struct TunnelConfig {
    pub private_key: String,
    pub addresses: Vec<String>,
    pub dns_servers: Vec<String>,
    pub mtu: u16,
    pub peer: PeerConfig,
}

impl TunnelConfig {
    pub fn validate(&self) -> Result<(), &'static str> {
        if self.private_key.trim().is_empty() {
            return Err("private key is required");
        }
        if self.addresses.is_empty() {
            return Err("at least one tunnel address is required");
        }
        if self.peer.public_key.trim().is_empty() {
            return Err("peer public key is required");
        }
        if self.peer.endpoint.trim().is_empty() {
            return Err("peer endpoint is required");
        }
        if self.peer.allowed_ips.is_empty() {
            return Err("at least one allowed IP is required");
        }
        if !(576..=9000).contains(&self.mtu) {
            return Err("mtu outside supported range");
        }
        Ok(())
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    fn valid() -> TunnelConfig {
        TunnelConfig {
            private_key: "runtime-secret-not-a-real-key".into(),
            addresses: vec!["10.88.0.2/32".into()],
            dns_servers: vec!["1.1.1.1".into()],
            mtu: 1420,
            peer: PeerConfig {
                public_key: "runtime-peer-key-not-a-real-key".into(),
                endpoint: "vpn.example.invalid:51820".into(),
                allowed_ips: vec!["0.0.0.0/0".into()],
                persistent_keepalive_seconds: Some(25),
            },
        }
    }

    #[test]
    fn valid_shape_passes() {
        assert!(valid().validate().is_ok());
    }

    #[test]
    fn missing_private_key_fails() {
        let mut config = valid();
        config.private_key.clear();
        assert!(config.validate().is_err());
    }

    #[test]
    fn unreasonable_mtu_fails() {
        let mut config = valid();
        config.mtu = 100;
        assert!(config.validate().is_err());
    }
}
