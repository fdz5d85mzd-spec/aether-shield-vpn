use serde::{Deserialize,Serialize};
#[derive(Debug,Clone,Serialize,Deserialize,PartialEq)]
pub struct Telemetry { pub rtt_ms:u32, pub packet_loss_pct:f32, pub transport_available:bool, pub consecutive_failures:u8 }
