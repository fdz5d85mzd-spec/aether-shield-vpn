package sessions

import (
 "errors"
 "github.com/fdz5d85mzd-spec/aether-shield-vpn/aether-shield-server/internal/nodes"
)

var ErrNoNode = errors.New("no verified online node available")

type Request struct { Region string; ClientPublicKey string }
type Provision struct { NodeID string `json:"nodeId"`; Endpoint string `json:"endpoint"`; ServerPublicKey string `json:"serverPublicKey"` }

type Provisioner struct{ registry *nodes.Registry }
func NewProvisioner(r *nodes.Registry)*Provisioner{return &Provisioner{registry:r}}

func (p *Provisioner) Provision(req Request)(Provision,error){
 if req.ClientPublicKey==""{return Provision{},errors.New("client public key required")}
 n,ok:=p.registry.SelectOnline(req.Region);if !ok{return Provision{},ErrNoNode}
 return Provision{NodeID:n.ID,Endpoint:n.Endpoint,ServerPublicKey:n.PublicKey},nil
}
