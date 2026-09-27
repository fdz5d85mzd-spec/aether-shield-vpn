pub mod android;
pub mod ffi;
pub mod policy;
pub mod telemetry;
pub mod tunnel;

#[no_mangle]
pub extern "C" fn aether_core_init() -> i32 {
    0
}
