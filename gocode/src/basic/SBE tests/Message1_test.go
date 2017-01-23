package sbe_tests

import (
	"bytes"
	"encoding/binary"
	"testing"
)

func TestEncodeDecodeMessage1(t *testing.T) {

	var in Message1
	copy(in.EDTField[:], "abcdefghijklmnopqrst")
	in.Header = MessageHeader{in.SbeBlockLength(), in.SbeTemplateId(), in.SbeSchemaId(), in.SbeSchemaVersion()}
	in.ENUMField = ENUM.Value10

	var buf = new(bytes.Buffer)
	if err := in.Encode(buf, binary.LittleEndian, true); err != nil {
		t.Logf("Encoding Error", err)
		t.Fail()
	}

	var out Message1 = *new(Message1)
	if err := out.Decode(buf, binary.LittleEndian, in.SbeSchemaVersion(), in.SbeBlockLength(), true); err != nil {
		t.Logf("Decoding Error", err)
		t.Fail()
	}

	if in != out {
		t.Logf("in != out:\n", in, out)
		t.Fail()
	}

	// Check that if we pass in an unknown enum value we get an error
	in.ENUMField = 77
	if err := in.Encode(buf, binary.LittleEndian, true); err == nil {
		t.Logf("Encoding didn't fail on unknown value with rangecheck")
		t.Fail()
	} else {
		t.Logf("Error looks good:", err)
	}

	// Encode should work if we don't rangecheck
	buf = new(bytes.Buffer)
	if err := in.Encode(buf, binary.LittleEndian, false); err != nil {
		t.Logf("Encoding failed on unknown value without rangecheck")
		t.Fail()
	}

	// Decode should fail on the bogus value if we rangecheck
	// FIXME: I reckon this should work
	if err := out.Decode(buf, binary.LittleEndian, in.SbeSchemaVersion(), in.SbeBlockLength(), true); err == nil {
		t.Logf("Decoding didn't fail on unknown value with rangecheck")
		t.Fail()
	} else {
		t.Logf("Error looks good:", err)
	}

	// Now check with proper values
	in.ENUMField = ENUM.Value10
	in.Encode(buf, binary.LittleEndian, true)
	if err := out.Decode(buf, binary.LittleEndian, in.SbeSchemaVersion(), in.SbeBlockLength(), true); err != nil {
		t.Logf("Decoding Error:", err)
		t.Fail()
	}
	if out.ENUMField != in.ENUMField {
		t.Logf("out.ENUMField != in.ENUMField", out.ENUMField, in.ENUMField)
		t.Fail()
	}

	return
}