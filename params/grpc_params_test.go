package params

import "testing"

type GrpcParams struct {
	addr string // grpc server address
	code int    // grpc code
	msg  string // grpc message
}

type GrpcParamsOption func(*GrpcParams) // GrpcParamsOption grpc params option

func NewAddr(addr string) GrpcParamsOption {
	return func(p *GrpcParams) {
		p.addr = addr
	}
}

func NewCode(code int) GrpcParamsOption {
	return func(p *GrpcParams) {
		p.code = code
	}
}

func NewMsg(msg string) GrpcParamsOption {
	return func(p *GrpcParams) {
		p.msg = msg
	}
}

func NewGrpcParams(opts ...GrpcParamsOption) *GrpcParams {
	p := &GrpcParams{
		addr: "localhost:80", // 默认地址
		code: 0,              // 默认code
		msg:  "",             // 默认消息
	}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

func TestCase1(t *testing.T) {
	p := NewGrpcParams(
		NewAddr("localhost:50051"),
		NewCode(200),
		NewMsg("OK"),
	)
	if p.addr != "localhost:50051" || p.code != 200 || p.msg != "OK" {
		t.Errorf("TestCase1 failed, got %+v", p)
	} else {
		t.Logf("TestCase1 passed, got %+v", p)
	}
}
func TestCase2(t *testing.T) {
	p := NewGrpcParams(
		NewAddr("localhost:50051"),
		NewMsg("OK"),
	)
	if p.code == 200 {
		t.Errorf("TestCase1 failed, got %+v", p)
	} else {
		t.Logf("TestCase1 passed, got %+v", p)
	}
}
