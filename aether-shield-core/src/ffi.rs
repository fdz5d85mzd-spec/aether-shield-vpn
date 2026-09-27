use std::sync::{Mutex, OnceLock};

use crate::session::TunnelSession;
use crate::transport::UnconfiguredTransport;

fn session() -> &'static Mutex<TunnelSession<UnconfiguredTransport>> {
    static INSTANCE: OnceLock<Mutex<TunnelSession<UnconfiguredTransport>>> = OnceLock::new();
    INSTANCE.get_or_init(|| Mutex::new(TunnelSession::new(UnconfiguredTransport)))
}

/// Attaches the platform TUN descriptor to the native core.
///
/// Return values:
/// 0 = TUN accepted; transport is NOT necessarily connected.
/// -2 = invalid state or descriptor.
/// -3 = internal synchronization error.
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

/// Attempts to establish packet transport after the TUN is attached.
///
/// The current build deliberately returns -4 because no production transport
/// has been configured yet. This prevents native clients from claiming a
/// connected VPN merely because the OS TUN interface exists.
#[no_mangle]
pub extern "C" fn aether_transport_connect() -> i32 {
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
    match session().lock() {
        Ok(mut session) => {
            session.stop();
            0
        }
        Err(_) => -3,
    }
}
