package example

import "context"

type Ret struct {
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
	Code int         `json:"code"`
}

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
		String  string
		Int     int
		Float32 float32
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
	// @test:http.delete(route="/:add2")
	//
	// @test2:http.get("/add3","/add4")
	//
	// ignore this
	// ignore this
	Add2(ctx context.Context, param Param) (res Ret, err error)
	// dec func
	// Doc for 1
	// @http.post(method=post,route="/dec")
	Dec(ctx context.Context, param Param) (err error)
	// dec func
	// Doc for 1
	// @http.get(method=post,route="/:dec")
	Dec2(ctx context.Context, param Param)

	// ignore this
	// @test2:http.post(method=get,route="/dec")
	TestService
}
