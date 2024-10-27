package dump

import (
	"testing"

	rtest "github.com/hantbk/vbackup/pkg/restic_source/rinternal/test"
)

func TestFormatLinuxACL(t *testing.T) {
	for _, c := range []struct {
		in, out, err string
	}{
		{
			in: "\x02\x00\x00\x00\x01\x00\x06\x00\xff\xff\xff\xff\x02\x00" +
				"\x04\x00\x03\x00\x00\x00\x02\x00\x04\x00\xe9\x03\x00\x00" +
				"\x04\x00\x02\x00\xff\xff\xff\xff\b\x00\x01\x00'\x00\x00\x00" +
				"\x10\x00\a\x00\xff\xff\xff\xff \x00\x04\x00\xff\xff\xff\xff",
			out: "user::rw-\nuser:3:r--\nuser:1001:r--\ngroup::-w-\n" +
				"group:39:--x\nmask::rwx\nother::r--\n",
		},
		{
			in: "\x02\x00\x00\x00\x00\x00\x06\x00\xff\xff\xff\xff\x02\x00" +
				"\x04\x00\x03\x00\x00\x00\x02\x00\x04\x00\xe9\x03\x00\x00" +
				"\x04\x00\x06\x00\xff\xff\xff\xff\b\x00\x05\x00'\x00\x00\x00" +
				"\x10\x00\a\x00\xff\xff\xff\xff \x00\x04\x00\xff\xff\xff\xff",
			err: "unknown tag",
		},
		{
			in: "\x01\x00\x00\x00\x01\x00\x06\x00\xff\xff\xff\xff\x02\x00" +
				"\x04\x00\x03\x00\x00\x00\x02\x00\x04\x00\xe9\x03\x00\x00" +
				"\x04\x00\x06\x00\xff\xff\xff\xff\b\x00\x05\x00'\x00\x00\x00" +
				"\x10\x00\a\x00\xff\xff\xff\xff \x00\x04\x00\xff\xff\xff\xff",
			err: "unsupported ACL format version",
		},
		{in: "\x02\x00", err: "wrong length"},
		{in: "", err: "wrong length"},
	} {
		out, err := formatLinuxACL([]byte(c.in))
		if c.err == "" {
			rtest.Equals(t, c.out, out)
		} else {
			rtest.Assert(t, err != nil, "wanted %q but got nil", c.err)
			rtest.Equals(t, c.err, err.Error())
		}
	}
}
