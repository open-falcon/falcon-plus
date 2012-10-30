#include <stdlib.h>
#include <rrd.h>

void rrdGetContext() {
	rrd_context_t *ctx = rrd_new_context();
	if (ctx == NULL) {
		//runtime·throw("librrd: out of memory");
	}
}

char *rrdCreate(const char *filename, unsigned long step, time_t start, int argc, const char **argv) {
	rrdGetContext();
	rrd_create_r(filename, step, start, argc, argv);
	if (rrd_test_error()) {
		return rrd_get_error();
	}
	return NULL;
}
