use std::sync::{Mutex, OnceLock};
use crate::tunnel::TunnelLifecycle;

fn lifecycle() -> &'static Mutex<TunnelLifecycle> {
    static INSTANCE: OnceLock<Mutex<TunnelLifecycle>> = OnceLock::new();
    INSTANCE.get_or_init(|| Mutex::new(TunnelLifecycle::default()))
}

#[no_mangle]
pub extern "C" fn aether_tunnel_start(tun_fd: i32) -> i32 {
    match lifecycle().lock() {
        Ok(mut tunnel) => match tunnel.start(tun_fd) { Ok(()) => 0, Err(_) => -2 },
        Err(_) => -3,
    }
}

#[no_mangle]
pub extern "C" fn aether_tunnel_stop() -> i32 {
    match lifecycle().lock() {
        Ok(mut tunnel) => { tunnel.stop(); 0 }
        Err(_) => -3,
    }
}
