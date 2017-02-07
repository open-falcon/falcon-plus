import json
import socket
import itertools
import time
 
class RPCClient(object):
 
    def __init__(self, addr, codec=json):
        self._socket = socket.create_connection(addr)
        self._id_iter = itertools.count()
        self._codec = codec
 
    def _message(self, name, *params):
        return dict(id=self._id_iter.next(),
                    params=list(params),
                    method=name)
 
    def call(self, name, *params):
        req = self._message(name, *params)
        id = req.get('id')
 
        mesg = self._codec.dumps(req)
        self._socket.sendall(mesg)
 
        # This will actually have to loop if resp is bigger
        resp = self._socket.recv(4096)
        resp = self._codec.loads(resp)
 
        if resp.get('id') != id:
            raise Exception("expected id=%s, received id=%s: %s"
                            %(id, resp.get('id'), resp.get('error')))
 
        if resp.get('error') is not None:
            raise Exception(resp.get('error'))
 
        return resp.get('result')
 
    def close(self):
        self._socket.close()
 

if __name__ == '__main__':
    rpc = RPCClient(("127.0.0.1", 8433))
    for i in xrange(10000):
        mv1 = dict(endpoint='host.niean', metric='metric.niean.1', value=i, step=60, 
            counterType='GAUGE', tags='tag=t'+str(i), timestamp=int(time.time()))
        mv2 = dict(endpoint='host.niean', metric='metric.niean.2', value=i, step=60, 
            counterType='COUNTER', tags='tag=t'+str(i), timestamp=int(time.time()))
        print rpc.call("Transfer.Update", [mv1, mv2])
