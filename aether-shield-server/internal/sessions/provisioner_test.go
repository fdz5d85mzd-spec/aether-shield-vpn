package sessions
import("errors";"testing";"github.com/fdz5d85mzd-spec/aether-shield-vpn/aether-shield-server/internal/nodes")
func TestProvisionFailsClosedWithoutInfrastructure(t *testing.T){p:=NewProvisioner(nodes.NewRegistry());_,err:=p.Provision(Request{ClientPublicKey:"client-public"});if !errors.Is(err,ErrNoNode){t.Fatalf("got %v",err)}}
func TestProvisionReturnsOnlyPublicNodeMaterial(t *testing.T){r:=nodes.NewRegistry();_ = r.Upsert(nodes.Node{ID:"n1",Region:"eu",Endpoint:"example.invalid:51820",PublicKey:"server-public",Status:nodes.StatusOnline});p:=NewProvisioner(r);v,err:=p.Provision(Request{Region:"eu",ClientPublicKey:"client-public"});if err!=nil{t.Fatal(err)};if v.ServerPublicKey!="server-public"{t.Fatal("wrong public key")}}
