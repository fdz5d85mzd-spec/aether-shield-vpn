package nodes
import "testing"
func TestRegistryDoesNotInventNodes(t *testing.T){r:=NewRegistry();if got:=r.List();len(got)!=0{t.Fatalf("got %d nodes",len(got))}}
func TestSelectRequiresOnline(t *testing.T){r:=NewRegistry();_ = r.Upsert(Node{ID:"n1",Region:"eu",Endpoint:"example.invalid:51820",PublicKey:"public-only",Status:StatusUnknown});if _,ok:=r.SelectOnline("eu");ok{t.Fatal("unknown node selected as online")}}
func TestUpsertRejectsMissingPublicData(t *testing.T){r:=NewRegistry();if err:=r.Upsert(Node{ID:"n1"});err==nil{t.Fatal("expected validation error")}}
