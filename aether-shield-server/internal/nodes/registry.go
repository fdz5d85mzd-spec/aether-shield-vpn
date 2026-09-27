package nodes

import (
 "errors"
 "sync"
 "time"
)

type Status string
const (
 StatusUnknown Status = "UNKNOWN"
 StatusOnline Status = "ONLINE"
 StatusOffline Status = "OFFLINE"
)

type Node struct {
 ID string `json:"id"`
 Region string `json:"region"`
 City string `json:"city"`
 Country string `json:"country"`
 Endpoint string `json:"endpoint"`
 PublicKey string `json:"publicKey"`
 Status Status `json:"status"`
 LastSeen *time.Time `json:"lastSeen"`
}

type Registry struct{ mu sync.RWMutex; nodes map[string]Node }

func NewRegistry() *Registry { return &Registry{nodes: make(map[string]Node)} }

func (r *Registry) Upsert(n Node) error {
 if n.ID=="" || n.Endpoint=="" || n.PublicKey=="" { return errors.New("id, endpoint and public key are required") }
 r.mu.Lock(); defer r.mu.Unlock(); r.nodes[n.ID]=n; return nil
}

func (r *Registry) Get(id string)(Node,bool){ r.mu.RLock(); defer r.mu.RUnlock(); n,ok:=r.nodes[id]; return n,ok }

func (r *Registry) List()[]Node{
 r.mu.RLock(); defer r.mu.RUnlock()
 out:=make([]Node,0,len(r.nodes)); for _,n:=range r.nodes{out=append(out,n)}; return out
}

func (r *Registry) SelectOnline(region string)(Node,bool){
 r.mu.RLock(); defer r.mu.RUnlock()
 for _,n:=range r.nodes{if n.Status==StatusOnline && (region=="" || n.Region==region){return n,true}}
 return Node{},false
}
