use crate::config::TunnelConfig;
use crate::engine::{EngineState, EngineTelemetry, PacketEngine};

pub struct EngineRuntime<E: PacketEngine> {
    engine: E,
    config: Option<TunnelConfig>,
    tun_fd: Option<i32>,
}

impl<E: PacketEngine> EngineRuntime<E> {
    pub fn new(engine: E) -> Self {
        Self {
            engine,
            config: None,
            tun_fd: None,
        }
    }

    pub fn configure(&mut self, config: TunnelConfig) -> Result<(), &'static str> {
        config.validate()?;
        self.engine.configure(&config)?;
        self.config = Some(config);
        Ok(())
    }

    pub fn attach_tun(&mut self, tun_fd: i32) -> Result<(), &'static str> {
        if tun_fd < 0 {
            return Err("invalid tun fd");
        }
        self.tun_fd = Some(tun_fd);
        Ok(())
    }

    pub fn start(&mut self) -> Result<(), &'static str> {
        if self.config.is_none() {
            return Err("runtime not configured");
        }
        let fd = self.tun_fd.ok_or("tun not attached")?;
        self.engine.start(fd)?;
        if self.engine.state() != EngineState::Active {
            return Err("packet engine not active");
        }
        Ok(())
    }

    pub fn stop(&mut self) {
        self.engine.stop();
        self.tun_fd = None;
        self.config = None;
    }

    pub fn telemetry(&self) -> EngineTelemetry {
        self.engine.telemetry()
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::engine::UnsupportedEngine;

    #[test]
    fn unsupported_runtime_fails_before_forwarding() {
        let mut runtime = EngineRuntime::new(UnsupportedEngine);
        assert!(runtime.attach_tun(8).is_ok());
        assert!(runtime.start().is_err());
        assert_eq!(runtime.telemetry(), EngineTelemetry::default());
    }
}
