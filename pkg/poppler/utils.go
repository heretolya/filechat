package poppler

// #include <glib.h>
// #include <unistd.h>
// #include <stdlib.h>
import "C"
import "unsafe"

func toString(in *C.gchar) string {
	str := C.GoString((*C.char)(in))
	C.free(unsafe.Pointer(in))
	return str
}

func toBool(in C.gboolean) bool {
	return int(in) > 0
}
