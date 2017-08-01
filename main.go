package main

import (
	"github.com/zhutingle/gotrix/server"
)

func main() {

	server.GotrixServer()

}

//func testSignal() {
//	signs := make(chan os.Signal, 1)
//	done := make(chan bool, 1)
//
//	signal.Notify(signs, syscall.SIGINT, syscall.SIGUSR1, syscall.SIGUSR2, syscall.SIGHUP, syscall.SIGTERM)
//
//	fmt.Println(os.Getpid())
//
//	go func() {
//		for {
//			sig := <-signs
//			fmt.Println(sig)
//			done <- true
//		}
//
//	}()
//
//	for i := 0; i < 10; i++ {
//		<-done
//		fmt.Println("Done Signal")
//	}
//}

//func checkError(err error) {
//	if err != nil {
//		fmt.Fprintf(os.Stderr, "Fatal error: %v", err)
//		os.Exit(1)
//	}
//}
//
//func testPrint() {
//
//	tcpAddr, err := net.ResolveTCPAddr("tcp4", "192.168.0.201:6101")
//	checkError(err)
//
//	conn, err := net.DialTCP("tcp", nil, tcpAddr)
//	checkError(err)
//
//	strs := "CT~~CD,~CC^~CT~\n";
//	strs += "^XA~TA000~JSN^LT0^MNW^MTD^PON^PMN^LH0,0^JMA^PR4,4~SD15^JUS^LRN^CI0^XZ\n"
//	strs += "^XA\n"
//	strs += "^MMT\n"
//	strs += "^PW827\n"
//	strs += "^LL0413\n"
//	strs += "^LS0\n"
//	strs += "^FO0,0^GFA,43264,43264,00104,:Z64:\n"
//	strs += "eJzs3U1uwjAQBlBXLLrkCByFo5Wj5Sg5QpcsUKcKjp2koiuCUaL3rawJ8luO/JOQkoiIiIjsJbHITy6eo0/HuKbD8umVw+FwOJzNOZc6PhRnSHb+FDgcDofD2Z5zHrvlfdpc6Cc499Kew+FwOJw9OCXHCeZwOBwOZ0fOV3TpFMHhcDgczj6ckuoM+U4PNkw5HA6Hw9maM89sf/TV92o4HA6Hw2ngPICnc8BXhsPhcDicBk7ulEP7rCn9dGywQ56HORwOh8N5j3MZR59xSx8R0/5oLeSstQ/L4XA4HE5LJ+q0Mb/vkguR4e4JgcPhcDic9zlllmGZeMtwuegZuZ8Oyb947hyQw+FwOJz2Tpk2X+OMvEzMhdpPh5zi+e+LcTgcDofT2vkH7u/jw3IhuWo4HA6Hw2ngxCLj6+ddmi0Ty0nh6u/pcTgcDofTwLnUcf2O5qKQc4qV/xeBw+FwOJxGzvneS2efE+tSSsdHBQ6Hw+FwNu7M73XWFxrW2B9d3/kFAAD//+3bMQ4AEBBFQfc/tQIbdCJLM+8CU25CPofD4XA4x047n8t+ofSBesLd5nA4HA4n3RnF+2iDt/3CXRwOh8Ph/HHm4h8w4Z5yOBwOh/PekSRJkiSpAtRO3Rs=:7163\n"
//	strs += "^PQ1,0,1,Y^XZ\n"
//
//	_, err = conn.Write([]byte(strs))
//	checkError(err)
//
//	result, err := ioutil.ReadAll(conn)
//	checkError(err)
//	fmt.Println(string(result))
//	os.Exit(0)
//}
