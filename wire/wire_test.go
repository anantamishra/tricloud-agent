package wire

import (
	"fmt"
	"testing"
)

func TestHeader(t *testing.T) {

	mybody := []byte("hellomessage")

	connid := UID(12)
	cmdtype := CommandType(93)

	packet := AttachHeader(connid, cmdtype, UserToServer, mybody)
	header := GetHeader(packet)

	t.Logf("%d", packet)
	t.Logf("%+v", header)

	if header.Connid != connid {
		t.Error("worong connid")
	}

	if header.CmdType != cmdtype {
		t.Error("worong msgtype")
	}
	if header.Flow != UserToServer {
		t.Error("worong flow")
	}
}

func TestEncodingDecoding(t *testing.T) {
	s := &SysStatData{
		AvailableMem: 200,
		TotalMem:     500,
	}

	bite, err := Encode(10, CMD_SYSTEMSTAT, AgentToServer, s)
	if err != nil {
		t.Error(err)
	}

	s2 := &SysStatData{}

	h, err := Decode(bite, s2)
	if err != nil {
		t.Error(err)
	}

	t.Logf("%+v", h)
	t.Logf("%+v", s2)

	assertEqual(t, s.AvailableMem, s2.AvailableMem, "")
	assertEqual(t, s.TotalMem, s2.TotalMem, "")

}

func assertEqual(t *testing.T, a interface{}, b interface{}, message string) {
	if a == b {
		return
	}
	if len(message) == 0 {
		message = fmt.Sprintf("%v != %v", a, b)
	}
	t.Fatal(message)
}
