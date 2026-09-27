pub mod android;
pub mod config;
pub mod ffi;
pub mod keys;
pub mod policy;
pub mod session;
pub mod telemetry;
pub mod transport;
pub mod tunnel;
pub mod wireguard;

#[no_mangle]
pub extern "C" fn aether_core_init() -> i32 {
    0
}
