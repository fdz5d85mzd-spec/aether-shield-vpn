use crate::engine::{EngineState, EngineTelemetry};

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub struct RuntimeStatus {
    pub engine_state: EngineState,
    pub telemetry: EngineTelemetry,
}

impl RuntimeStatus {
    pub fn is_connected(&self) -> bool {
        self.engine_state == EngineState::Active
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn only_active_engine_is_connected() {
        let empty = EngineTelemetry::default();
        assert!(
            !RuntimeStatus {
                engine_state: EngineState::Unsupported,
                telemetry: empty,
            }
            .is_connected()
        );
        assert!(
            !RuntimeStatus {
                engine_state: EngineState::Ready,
                telemetry: empty,
            }
            .is_connected()
        );
        assert!(
            RuntimeStatus {
                engine_state: EngineState::Active,
                telemetry: empty,
            }
            .is_connected()
        );
    }
}
