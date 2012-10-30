extern char *rrdCreate(const char *filename, unsigned long step, time_t start, int argc, const char **argv);
extern char *rrdUpdate(const char *filename, const char *template, int argc, const char **argv);
