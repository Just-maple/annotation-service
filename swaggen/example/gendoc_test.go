package example

import (
	"os/exec"
	"testing"

	"github.com/Just-maple/annotation-service/swaggen"
)

func TestGen(t *testing.T) {
	swaggen.GenDoc("./", "./define.go",
		swaggen.WithRetPack(&Ret{}, "data"),
	)
	err := exec.Command("swag", "init", "--generalInfo", "define.go").Run()
	if err != nil {
		panic(err)
	}
}
