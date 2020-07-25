package annotation_service

import (
	"context"
	"runtime"
	"testing"
)

type (
	// @service(test)
	TestService interface {
		// add func
		// @http(method=get,route="/add")
		Add(ctx context.Context, param Param) (res int, err error)
	}

	// @service(test)
	TestService3 interface {
		// add func
		// @http(method=get,route="/add")
		Add(ctx context.Context, param Param) (res int, err error)
	}

	Param struct {
	}
)

// @service(test)
type TestService2 interface {
	// add func
	// @test:http.get(method=get,route="/add")
	Add2(ctx context.Context, param Param) (res int, err error)
	// dec func
	// @http.get.post(method=get,route="/dec")
	Dec(ctx context.Context, param Param) (res int, err error)

	// ignore this
	// @http.get.post(method=get,route="/dec")
	TestService
}

func TestGetAllService(t *testing.T) {
	_, f, _, _ := runtime.Caller(0)
	ret, err := GetAllService(f)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%v", ret)
}
