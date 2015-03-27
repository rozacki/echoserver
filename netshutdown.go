package main

import ("net"
 "sync"
"log")

type gListener struct{
	net.Listener
	tcpWA	*sync.WaitGroup
	stop chan error
	closed bool
}

func newgListener(addr string) (gL* gListener, err error){

	//create listener that listens on all interfaces on port...
	l,err:=net.Listen("tcp4",addr)
	if err!=nil{
		return
	}

	gL=&gListener{l, &sync.WaitGroup{},make(chan error),false}

	//wait for signal and close listener
	go func(){
		//wait for signal to stop
		_=<-gL.stop
		log.Println("signal to stop")
		err:=gL.Listener.Close()
		gL.closed=true
		//send back that closing is finished
		gL.stop<-err
	}()

	return
}

func (self *gListener) Accept() (gC* gConn, err error){
	conn,err:= self.Listener.Accept()

	if err!=nil{
		return nil, err
	}

	log.Println("added to wait group")
	self.tcpWA.Add(1)

	gC=&gConn{ Conn: conn, tcpWA:self.tcpWA}

	return
}

func (self* gListener) Close() error{
	if self.closed==true{
		return nil
	}
	return self.Listener.Close()
}

type gConn struct{
	net.Conn
	tcpWA* sync.WaitGroup
}

func (self gConn) Close() error {

	log.Println("wait group done")
	self.tcpWA.Done()

	return self.Conn.Close()
}
