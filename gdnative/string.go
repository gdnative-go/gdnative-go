package gdnative

/*
#include <gdnative/string.h>
#include "gdnative.gen.h"
*/
import "C"

// NewStringWithWideString creates a new String with given contents
func NewStringWithWideString(str string) String {
	return String(str)
}

// NewString retruns an empty String
func NewString() String {
	return ""
}

func truncateString(str string, num int) string {
	return str[0:num]
}
