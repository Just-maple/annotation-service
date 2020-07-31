package mount_field

import (
	"testing"

	"github.com/Just-maple/annotation-service/mount_field/testfield"
)

type Test struct {
	A testfield.A
	B testfield.B
	//comment1
	//comment2
}

func TestMount(t *testing.T) {
	set, err := NewStructMounter("./field_mount_test.go", "Test")
	if err != nil {
		t.Fatal(err)
	}

	_ = set.MountTypeField("testfield.A")
	_ = set.MountTypeField("testfield.B")

	err = set.Write()
	if err != nil {
		t.Fatal(err)
	}
	return
}
