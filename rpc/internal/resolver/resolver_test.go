package resolver

import (
	"testing"

	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/serviceconfig"
)

func TestNopResolver(t *testing.T) {
	//	make sure ResolveNow & Close don't panic
	var r nopResolver
	r.ResolveNow(resolver.ResolveNowOptions{})
	r.Close()
}

type mockedClientConn struct {
	state resolver.State
}

func (m *mockedClientConn) ReportError(err error) {
	panic("implement me")
}

func (m *mockedClientConn) NewAddress(addresses []resolver.Address) {
	panic("implement me")
}

func (m *mockedClientConn) NewServiceConfig(serviceConfig string) {
	panic("implement me")
}

func (m *mockedClientConn) ParseServiceConfig(serviceConfigJSON string) *serviceconfig.ParseResult {
	panic("implement me")
}

func (m *mockedClientConn) UpdateState(state resolver.State) error {
	m.state = state
	return nil
}
