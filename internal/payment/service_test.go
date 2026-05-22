package payment

import "testing"

func TestNewService(t *testing.T) {
	if New(nil, nil) == nil {
		t.Fatal("service should be constructed")
	}
}
