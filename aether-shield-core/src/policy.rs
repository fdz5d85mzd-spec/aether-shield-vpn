use crate::telemetry::Telemetry;
#[derive(Debug,Clone,Copy,PartialEq,Eq)]
pub enum TransportPolicy { WireGuardUdp, ResilientTransportRequired, Offline }
pub fn evaluate(t:&Telemetry)->TransportPolicy{
 if !t.transport_available || t.consecutive_failures>=3 { return TransportPolicy::Offline; }
 if t.packet_loss_pct>4.0 || t.rtt_ms>100 { return TransportPolicy::ResilientTransportRequired; }
 TransportPolicy::WireGuardUdp
}
#[cfg(test)]
mod tests{
 use super::*; use crate::telemetry::Telemetry;
 #[test] fn healthy_prefers_wireguard(){assert_eq!(evaluate(&Telemetry{rtt_ms:25,packet_loss_pct:0.2,transport_available:true,consecutive_failures:0}),TransportPolicy::WireGuardUdp)}
 #[test] fn degraded_requires_resilient_transport(){assert_eq!(evaluate(&Telemetry{rtt_ms:140,packet_loss_pct:1.0,transport_available:true,consecutive_failures:0}),TransportPolicy::ResilientTransportRequired)}
 #[test] fn unavailable_is_offline(){assert_eq!(evaluate(&Telemetry{rtt_ms:0,packet_loss_pct:0.0,transport_available:false,consecutive_failures:0}),TransportPolicy::Offline)}
}
