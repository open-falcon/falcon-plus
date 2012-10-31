#include <stdlib.h>
#include <rrd.h>

char *rrdError() {
	char *err = NULL;
	if (rrd_test_error()) {
		// RRD error is local for thread so other gorutine can call some RRD
		// function in the same thread before we use C.GoString. So we need to
		// copy current error before return from C to Go. It need to be freed
		// after C.GoString in Go code.
		err = strdup(rrd_get_error());
		if (err == NULL) {
			abort();
		}
	}
	return err;
}

char *rrdCreate(const char *filename, unsigned long step, time_t start, int argc, const char **argv) {
	rrd_clear_error();
	rrd_create_r(filename, step, start, argc, argv);
	return rrdError();
}

char *rrdUpdate(const char *filename, const char *template, int argc, const char **argv) {
	rrd_clear_error();
	rrd_update_r(filename, template, argc, argv);
	return rrdError();
}

char *rrdGraph(rrd_info_t **ret, int argc, char **argv) {
	rrd_clear_error();
	*ret = rrd_graph_v(argc, argv);
	return rrdError();
}
