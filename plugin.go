package malvales

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
	rt "github.com/arnodel/golua/runtime"
)

type Plugin interface{
	Load(*rt.Runtime) rt.Value
}

type entryRPCClient struct{
	client *rpc.Client
}

func (e *entryRPCClient) Load(rtm *rt.Runtime) rt.Value {
	var resp *rt.Value
	err := e.client.Call("Plugin.Load", rtm, &resp)
	if err != nil {
		// TODO: return nil (or some value to indicate err)
		panic(err)
	}

	return *resp
}

type entryRPCServer struct{
	P Plugin
}

func (s *entryRPCServer) Load(rtm *rt.Runtime, resp *rt.Value) error {
	*resp = s.P.Load(rtm)
	return nil
}

type Entry struct{
	P Plugin
}

func (e *Entry) Server(b *plugin.MuxBroker) (interface{}, error) {
	return &entryRPCServer{P: e.P}, nil
}

func (e *Entry) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &entryRPCClient{client: c}, nil
}
