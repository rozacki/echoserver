# echoserver
Gracefull shutdown TCPListener

This is a demonstration of how a gracafule tcplistsner shutdown can be implemented. It will wait for any signal and trigger graceful shutdown
on the listener. It will wait for all open connection to finish before program quits. main.go does not have  net or sync imports.
