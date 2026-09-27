#[cfg(target_os = "android")]
use jni::objects::JObject;
#[cfg(target_os = "android")]
use jni::sys::jint;
#[cfg(target_os = "android")]
use jni::JNIEnv;

#[cfg(target_os = "android")]
#[no_mangle]
pub extern "system" fn Java_com_aethershield_vpn_AetherVpnService_nativeTunnelStart(
    _env: JNIEnv,
    _this: JObject,
    tun_fd: jint,
) -> jint {
    crate::ffi::aether_tunnel_start(tun_fd)
}

#[cfg(target_os = "android")]
#[no_mangle]
pub extern "system" fn Java_com_aethershield_vpn_AetherVpnService_nativeTunnelStop(
    _env: JNIEnv,
    _this: JObject,
) -> jint {
    crate::ffi::aether_tunnel_stop()
}
