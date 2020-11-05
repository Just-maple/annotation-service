package mount_field

import (
	"os"
	"testing"

	testfieldagsags "github.com/Just-maple/annotation-service/mount_field/testfield"
)

type Test struct {
	GSSG os.Signal
	FFF  os.Signal
	B    testfieldagsags.B
	//	asggasgas
}

// asgsagasg
func TestMount(t *testing.T) {
	set, err := NewStructMounter("./field_mount_test.go", "Test")
	if err != nil {
		t.Fatal(err)
	}

	_ = set.MountTypeField("os.Signal", "GSSG", "os")
	_ = set.MountTypeField("os.Signal", "FFF", "os")
	_ = set.MountTypeField("testfieldagsags.B", "", "github.com/Just-maple/annotation-service/mount_field/testfield")

	err = set.Write()
	if err != nil {
		t.Fatal(err)
	}
	return
}
