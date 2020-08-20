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

	// @service(test4)
	TestService4 interface {
		// add func
		// @http(method=get,route="/add")
		Add(ctx context.Context, param Param) (res int, err error)
	}

	// @service(test3)
	TestService3 interface {
		// add func
		// @http(method=get,route="/add")
		Add(ctx context.Context, param []Param) (res map[Param]map[string][]Param, err error)
	}

	Param struct {
	}

	// @dao(test)
	TestDAO interface {
		Add(ctx context.Context, param Param) (res int, err error)
	}
)

// @service(test)
// @service(test2)
type TestService2 interface {
	// This is title
	//
	// Doc for 1
	// @test:http.get(route="/add")
	//
	// Doc for 2
	// @test:http.get(route="/add2")
	//
	// @test2:http.get("/add3","/add4")
	//
	// ignore this
	// ignore this
	Add2(ctx context.Context, param Param) (res Service, err error)
	// dec func
	// Doc for 1
	// @http.get.post(method=get,route="/dec")
	Dec(ctx context.Context, param Param) (res int, err error)
	Dec2(ctx context.Context, param Param) (res int, err error)

	// ignore this
	// @test2:http.get.post(method=get,route="/dec")
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

func TestGetAllServiceIdent(t *testing.T) {
	_, f, _, _ := runtime.Caller(0)
	ret, err := GetAllService(f, WithIdent("dao"))
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%v", ret)
}
