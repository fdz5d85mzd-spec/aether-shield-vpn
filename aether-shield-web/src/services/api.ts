import type { NodeInfo,SystemStatus } from '../types';
const base=import.meta.env.VITE_AETHER_API_URL ?? 'http://localhost:8080';
async function get<T>(path:string):Promise<T>{const r=await fetch(base+path);if(!r.ok)throw new Error(`API ${r.status}`);return r.json() as Promise<T>}
export const api={status:()=>get<SystemStatus>('/v1/system/status'),nodes:()=>get<NodeInfo[]>('/v1/nodes')};
