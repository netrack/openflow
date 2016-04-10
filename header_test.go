package of

import (
	"testing"
)

func TestHeaderSet(t *testing.T) {
	var hdr header

	err := hdr.Set(VersionHeaderKey, uint8(4))
	if err != nil {
		t.Fatal("Failed to set version: ", err.Error())
	}

	if hdr.Version != 4 {
		t.Fatalf("Wrong header version: ", hdr.Version)
	}

	err = hdr.Set(TypeHeaderKey, TypeHello)
	if err != nil {
		t.Fatal("Failed to set type: ", err.Error())
	}

	if hdr.Type != TypeHello {
		t.Fatal("Wrong header type: ", hdr.Type)
	}

	err = hdr.Set(XIDHeaderKey, uint32(123456789))
	if err != nil {
		t.Fatal("Failed to set xid: ", err.Error())
	}

	if hdr.XID != 123456789 {
		t.Fatal("Wrong header xid: ", hdr.XID)
	}

	err = hdr.Set(13, "boom")
	if err == nil {
		t.Fatal("Faied to report error")
	}
}
