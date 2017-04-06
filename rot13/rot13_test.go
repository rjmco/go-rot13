package rot13

import (
	"bytes"
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

var cases = []struct {
	in, want string
}{
	{"Lbh penpxrq gur pbqr!", "You cracked the code!"},
	{"You cracked the code!", "Lbh penpxrq gur pbqr!"},
	{"Why did the chicken cross the road?", "Jul qvq gur puvpxra pebff gur ebnq?"},
	{"Gb trg gb gur bgure fvqr!", "To get to the other side!"},
}

func TestDecode(t *testing.T) {
	encodingTestHelper(t, Decode, outputChars, inputChars)
}

func TestEncode(t *testing.T) {
	encodingTestHelper(t, Encode, inputChars, outputChars)
}

func TestRead(t *testing.T) {
	for _, v := range cases {
		t.Logf("LOG: Testing case in='%v', want='%v'", v.in, v.want)
		str_reader := strings.NewReader(v.in)
		rot_reader := Reader{str_reader}
		rot_bytes := make([]byte, str_reader.Size())
		_, err := rot_reader.Read(rot_bytes)
		if err != nil {
			t.Errorf("ERROR: Read(%v) failed.", v.in)
		}
		test_rslt := fmt.Sprintf("in='%v', want='%v', got='%v'",
			v.in, v.want, string(rot_bytes))
		if v.want != string(rot_bytes) {
			t.Errorf("ERROR: Test case failed. %v", test_rslt)
		} else {
			t.Logf("LOG: Test case succeeded. %v", test_rslt)
		}
	}
}

func TestWrite(t *testing.T) {
	for _, v := range cases {
		t.Logf("LOG: Testing case in='%v', want='%v'", v.in, v.want)
		var s_buf bytes.Buffer
		rot_writer := Writer{&s_buf}
		_, err := rot_writer.Write([]byte(v.in))
		if err != nil {
			t.Errorf("ERROR: Write(%v) failed.", v.in)
		}
		got := s_buf.String()
		test_rslt := fmt.Sprintf("in='%v', want='%v', got='%v'",
			v.in, v.want, got)
		if got != v.want {
			t.Errorf("ERROR: Test case failed. %v", test_rslt)
		} else {
			t.Logf("LOG: Test case succeeded. %v", test_rslt)
		}

	}

}

func encodingTestHelper(t *testing.T, codeFunc func(byte) (byte, error), ic string, oc string) {
	t.Logf("LOG: encodingTestHelper invoked with:\n\tcodeFunc=%s\n\tic='%s'\n\toc='%s'",
		runtime.FuncForPC(reflect.ValueOf(codeFunc).Pointer()).Name(), ic, oc)
	//iterate over all possible byte sequences
	for b, first_pass := byte(0), true; first_pass; b++ {
		t.Logf("LOG: Testing byte %v", b)

		b_index := strings.IndexByte(ic, b)
		eb, err := codeFunc(byte(b))

		//if byte found in ic then encode
		if b_index != -1 {

			//make sure it is successfully encoded
			if err != nil {
				t.Errorf("ERROR: Failed to code %v", b)
			} else {
				t.Logf("LOG: %s was successfully coded", string(b))
			}

			//make sure it is correctly encoded
			if oc[b_index] != eb {
				t.Errorf("ERROR: Invalid coding of %s: got = %s, wanted = %s",
					string(b), string(eb), string(oc[b_index]))
			} else {
				t.Logf("LOG: %s was successfully coded into %s",
					string(b), string(oc[b_index]))
			}

		} else { // if byte is not found in ic

			//make sure encoding fails
			if err == nil {
				t.Errorf("ERROR: Unexpectedly coded invalid character %v", b)
			} else {
				t.Logf("LOG: EncodingError(): %s", err)
			}

			//make sure the right error type is passed
			switch err.(type) {
			case EncodingError:
				t.Logf("LOG: %v expectedly raised an EncodingError", b)
			default:
				t.Errorf("ERROR: Wrong error type %T, expected %T", err, EncodingError(0))
			}
		}

		//when b is about to overflow, stop the test
		if b == byte(1<<8-1) {
			t.Logf("LOG: Reached byte %v. Stopping test before overflowing", b)
			first_pass = false
		}
	}
}
