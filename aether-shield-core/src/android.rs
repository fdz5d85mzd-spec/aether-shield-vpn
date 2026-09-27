#[cfg(target_os = "android")]
use jni::objects::{JByteArray, JObject, JString};
#[cfg(target_os = "android")]
use jni::sys::{jbyteArray, jint, jstring};
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
pub extern "system" fn Java_com_aethershield_vpn_AetherVpnService_nativeTransportConnect(
    _env: JNIEnv,
    _this: JObject,
) -> jint {
    crate::ffi::aether_transport_connect()
}

#[cfg(target_os = "android")]
#[no_mangle]
pub extern "system" fn Java_com_aethershield_vpn_AetherVpnService_nativeTunnelStop(
    _env: JNIEnv,
    _this: JObject,
) -> jint {
    crate::ffi::aether_tunnel_stop()
}

#[cfg(target_os = "android")]
#[no_mangle]
pub extern "system" fn Java_com_aethershield_vpn_NativeKeys_nativeGeneratePrivateKey(
    env: JNIEnv,
    _this: JObject,
) -> jbyteArray {
    let pair = match crate::keys::KeyPair::generate() {
        Ok(pair) => pair,
        Err(_) => return std::ptr::null_mut(),
    };
    match env.byte_array_from_slice(&pair.private_bytes()) {
        Ok(array) => array.into_raw(),
        Err(_) => std::ptr::null_mut(),
    }
}

#[cfg(target_os = "android")]
#[no_mangle]
pub extern "system" fn Java_com_aethershield_vpn_NativeKeys_nativePublicKey(
    mut env: JNIEnv,
    _this: JObject,
    private_key: JByteArray,
) -> jstring {
    let bytes = match env.convert_byte_array(private_key) {
        Ok(bytes) if bytes.len() == 32 => bytes,
        _ => return std::ptr::null_mut(),
    };
    let mut key = [0u8; 32];
    key.copy_from_slice(&bytes);
    let pair = crate::keys::KeyPair::from_private_bytes(key);
    match env.new_string(pair.public_base64()) {
        Ok(value) => {
            let raw: JString = value;
            raw.into_raw()
        }
        Err(_) => std::ptr::null_mut(),
    }
}

#[cfg(target_os = "android")]
#[no_mangle]
pub extern "system" fn Java_com_aethershield_vpn_AetherVpnService_nativeConfigureTunnel(
    mut env: JNIEnv,
    _this: JObject,
    private_key: JByteArray,
    address: JString,
    endpoint: JString,
    server_public_key: JString,
) -> jint {
    use base64::{engine::general_purpose::STANDARD, Engine as _};

    let private = match env.convert_byte_array(private_key) {
        Ok(bytes) if bytes.len() == 32 => bytes,
        _ => return -10,
    };
    let address: String = match env.get_string(&address) {
        Ok(value) => value.into(),
        Err(_) => return -11,
    };
    let endpoint: String = match env.get_string(&endpoint) {
        Ok(value) => value.into(),
        Err(_) => return -11,
    };
    let server_public_key: String = match env.get_string(&server_public_key) {
        Ok(value) => value.into(),
        Err(_) => return -11,
    };

    let config = crate::ffi::build_runtime_config(
        STANDARD.encode(private),
        address,
        endpoint,
        server_public_key,
    );
    match crate::ffi::set_runtime_config(config) {
        Ok(()) => 0,
        Err(_) => -12,
    }
}
