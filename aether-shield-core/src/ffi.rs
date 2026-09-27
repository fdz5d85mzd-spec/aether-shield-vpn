use std::sync::{Mutex, OnceLock};

use crate::config::{PeerConfig, TunnelConfig};
use crate::session::TunnelSession;
use crate::transport::UnconfiguredTransport;

fn session() -> &'static Mutex<TunnelSession<UnconfiguredTransport>> {
    static INSTANCE: OnceLock<Mutex<TunnelSession<UnconfiguredTransport>>> = OnceLock::new();
    INSTANCE.get_or_init(|| Mutex::new(TunnelSession::new(UnconfiguredTransport)))
}

fn config_slot() -> &'static Mutex<Option<TunnelConfig>> {
    static CONFIG: OnceLock<Mutex<Option<TunnelConfig>>> = OnceLock::new();
    CONFIG.get_or_init(|| Mutex::new(None))
}

pub fn set_runtime_config(config: TunnelConfig) -> Result<(), &'static str> {
    config.validate()?;
    match config_slot().lock() {
        Ok(mut slot) => {
            *slot = Some(config);
            Ok(())
        }
        Err(_) => Err("configuration lock unavailable"),
    }
}

pub fn build_runtime_config(
    private_key: String,
    address: String,
    endpoint: String,
    server_public_key: String,
) -> TunnelConfig {
    TunnelConfig {
        private_key,
        addresses: vec![address],
        dns_servers: vec![],
        mtu: 1420,
        peer: PeerConfig {
            public_key: server_public_key,
            endpoint,
            allowed_ips: vec!["0.0.0.0/0".into()],
            persistent_keepalive_seconds: Some(25),
        },
    }
}

#[no_mangle]
pub extern "C" fn aether_tunnel_start(tun_fd: i32) -> i32 {
    match session().lock() {
        Ok(mut session) => match session.attach_tun(tun_fd) {
            Ok(()) => 0,
            Err(_) => -2,
        },
        Err(_) => -3,
    }
}

#[no_mangle]
pub extern "C" fn aether_transport_connect() -> i32 {
    if config_slot().lock().map_or(true, |slot| slot.is_none()) {
        return -5;
    }
    match session().lock() {
        Ok(mut session) => match session.connect() {
            Ok(()) => 0,
            Err(_) => -4,
        },
        Err(_) => -3,
    }
}

#[no_mangle]
pub extern "C" fn aether_tunnel_stop() -> i32 {
    if let Ok(mut slot) = config_slot().lock() {
        *slot = None;
    }
    match session().lock() {
        Ok(mut session) => {
            session.stop();
            0
        }
        Err(_) => -3,
    }
}
