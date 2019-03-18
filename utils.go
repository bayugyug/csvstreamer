package csvstreamer

import (
	"strings"
)

//fmtNumeric format the numeric into string
func fmtNumeric(s string) string {
	mlen := DecimalMaxLen
	t := s
	if len(t) > mlen {
		t = t[:mlen]
	}
	return t
}

//fmtDoubleQts format the double-qts
func fmtDoubleQts(s string) string {
	t := s
	t = strings.Replace(t, "\\", "", -1)
	t = strings.Replace(t, "\"", "", -1)
	t = strings.Replace(t, "\n", "", -1)
	t = strings.Replace(t, "\t", "", -1)
	t = strings.Replace(t, "\r", "", -1)
	t = strings.Replace(t, "'", "", -1)
	t = strings.Replace(t, "`", "", -1)
	t = strings.TrimSpace(fmtStripControlCharsAndExtUTF8(t, true))

	return t
}

//fmtStripControlCharsAndExtUTF8
//  stripped of control codes and extended characters
func fmtStripControlCharsAndExtUTF8(str string, utf8 bool) string {

	if !utf8 {
		//control-chars only
		return strings.Map(func(r rune) rune {
			if r >= 32 {
				return r
			}
			return -1
		}, str)
	}
	//control-chars and utf8 chars
	return strings.Map(func(r rune) rune {
		if r >= 32 && r < 127 {
			return r
		}
		return -1
	}, str)
}
