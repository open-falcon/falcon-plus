extern char *rrdCreate(const char *filename, unsigned long step, time_t start, int argc, const char **argv);
extern char *rrdUpdate(const char *filename, const char *template, int argc, const char **argv);
extern char *rrdGraph(rrd_info_t **ret, int argc, char **argv);
extern char *rrdInfo(rrd_info_t **ret, char *filename);
extern char *rrdFetch(int *ret, char *filename, const char *cf, time_t *start, time_t *end, unsigned long *step, unsigned long *ds_cnt, char ***ds_namv, double **data);
extern char *arrayGetCString(char **values, int i);
extern double arrayGetDouble(double *values, int i);
