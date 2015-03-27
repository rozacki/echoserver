package main

import(
	"log"
	"os/signal"
	"os"
)

func main(){

	gL,err:=newgListener("127.0.0.1:2000")

	if err!=nil{
		log.Fatal(err)
	}
	//when function exits then close listener
	defer gL.Close()

	ch:=make(chan os.Signal)
	signal.Notify(ch,)
	go func(chS chan os.Signal){
		for {
			//wait for signal
			s := <-chS
			log.Println("signal received: ", s)
			//notify listener to close
			gL.stop<-nil
			//leave gor
			return
		}
	}(ch)

	log.Println("listening on:",gL.Addr())

	//infinite loop of main goroutine that accepts all connections
	for {
		//if listener is closed then we quit
		if gL.closed==true {
			log.Println("listener close")
			break;
		}
		conn,err:=gL.Accept()
		if err!= nil{
			log.Println("accept error: ", err)
			continue
		}

		log.Printf("Connection accepted. Local address:%s, remote address:%s \n", conn.LocalAddr(), conn.RemoteAddr() )

		if err!=nil{
			log.Fatal(err)
		}

		//business logic
		go func(c* gConn){

			defer c.Close()

			buffer:=make([]byte,100)

			log.Println("before read")
			read,err:=c.Read(buffer)
			log.Println("after read")

			if err!=nil{
				log.Println("error while reading from socket:", err)
				return
			}

			log.Printf("read %d bytes,string read: %s\n",read, string(buffer))

			var written int
			written, err=c.Write(buffer)

			log.Printf("written %d bytes",written)

		}(conn)
	}
	//if I'm here it means main goroutine wants to quit
	//wait for all go routines to finish business in the main go routine
	gL.tcpWA.Wait()
}
