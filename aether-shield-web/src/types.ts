export type ServiceState='ONLINE'|'OFFLINE'|'UNKNOWN'|'DEMO';
export type AuditState='NOT_RUN'|'RUNNING'|'PASS'|'FAIL'|'ERROR'|'UNSUPPORTED';
export type EngineState='UNSUPPORTED'|'READY'|'ACTIVE'|'FAILED';

export interface NodeInfo {
  id:string;
  city:string;
  country:string;
  status:ServiceState;
  latencyMs:number|null;
  capabilities:string[];
  lastSeen:string|null;
}

export interface EngineTelemetry {
  packetsFromTun:number;
  packetsToTun:number;
  encryptedBytesSent:number;
  encryptedBytesReceived:number;
}

export interface SystemStatus {
  mode:'REAL'|'DEMO';
  tunnel:'NOT_CONFIGURED'|'DISCONNECTED'|'CONNECTING'|'CONNECTED'|'ERROR';
  api:string;
  engineState?:EngineState;
  telemetry?:EngineTelemetry;
}
