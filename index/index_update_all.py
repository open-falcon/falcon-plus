#!/bin/env python
# -*- coding: utf-8 -*-

import time, socket, json, httplib
import sys, getopt

urlUpdateAllFmt = "/index/updateAll?step=172800"

def update_index(host, port, url):
  result = None
  httpClient = None
  print "{host}:{port}{url}".format(host=host,port=port,url=url)
  try:
    httpClient = httplib.HTTPConnection(host, port, timeout=30)
    httpClient.request("GET", url)
    response = httpClient.getresponse()
    if response.status/100 == 2:
      result = json.loads(response.rad())["data"]
  except Exception, e:
    pass
  finally:
    if httpClient:
      httpClient.close()

  return result
  
def usage():
	print "-h hostname | -x"

if __name__ == "__main__":
	opts, args = getopt.getopt(sys.argv[1:], "xh:p:")
	host = ""
	port = ""
	for op, value in opts:
		if op == "-h":
			host = value
		if op == "-p":
			port = value
		elif op == "-x":
			usage()
			sys.exit()

	if len(host) <= 0 or len(port)<=0 :
		print "bad args"
		sys.exit()

	print update_index(host, port, urlUpdateAllFmt)	
