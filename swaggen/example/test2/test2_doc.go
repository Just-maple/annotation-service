package test2

import "github.com/Just-maple/annotation-service/swaggen/example"

// @Summary  This is title
// @Produce  json
// @Param query query example.Param true "example.Param"
// @Success 200 {object} example.Ret{data=example.Ret}
// @Failure 400 {object} example.Ret
// @Router /test2/add3 [get]
func _() { _ = example.TestService2.Add2 }

// @Summary  dec func
// @Produce  json
// @Param json body example.Param true "example.Param"
// @Accept json
// @Success 200 {object} example.Ret
// @Failure 400 {object} example.Ret
// @Router /test2/dec [post]
func _() { _ = example.TestService2.Dec }

// @Summary  dec func
// @Produce  json
// @Param query query example.Param true "example.Param"
// @Param dec path string true "dec"
// @Success 200 {object} example.Ret
// @Failure 400 {object} example.Ret
// @Router /test2/{dec} [get]
func _() { _ = example.TestService2.Dec2 }
