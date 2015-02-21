package ofp13

import (
	"bytes"
	"fmt"
	"testing"
)

func TestAction(t *testing.T) {
}

func TestActionOutput(t *testing.T) {
	var buf bytes.Buffer
	a := ActionOutput{P_FLOOD, 0}

	_, err := a.WriteTo(&buf)
	if err != nil {
		t.Fatal("Failed to marshal action:", err)
	}

	hexstr := fmt.Sprintf("%x", buf.Bytes())
	if hexstr != "00000010fffffffb0000000000000000" {
		t.Fatal("Marshaled action data is incorrect:", hexstr)
	}
}

func TestActionGroup(t *testing.T) {
}

func TestActionSetQueue(t *testing.T) {
}

func TestActionMPLSTTL(t *testing.T) {
}

func TestActionSetNetworkTTL(t *testing.T) {
}

func TestActionPush(t *testing.T) {
}

func TestActionPopMPLS(t *testing.T) {
}
