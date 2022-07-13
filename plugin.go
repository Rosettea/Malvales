package malvales

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
	rt "github.com/arnodel/golua/runtime"
)

type Module interface{
	Loader() rt.Value
}

type entryRPCClient struct{
	client *rpc.Client
}

func (e *entryRPCClient) Loader() rt.Value {
	var resp *rt.Value
	err := e.client.Call("Plugin.Loader", new(interface{}), &resp)
	if err != nil {
		// TODO: return nil (or some value to indicate err)
		panic(err)
	}

	return *resp
}

type entryRPCServer struct{
	M Module
}

func (s *entryRPCServer) Loader(_ interface{}, resp *rt.Value) error {
	*resp = s.M.Loader()
	return nil
}

type Entry struct{
	M Module
}

func (e *Entry) Server(b *plugin.MuxBroker) (interface{}, error) {
	return &entryRPCServer{M: e.M}, nil
}

func (e *Entry) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &entryRPCClient{client: c}, nil
}
