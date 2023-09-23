package main

import (
	"flag"
	"fmt"
	apizinc "github.com/JoseM-18/zincSearch/apiZinc"
	"github.com/JoseM-18/zincSearch/email"
	"github.com/JoseM-18/zincSearch/emailProcessing/finders"
	processor "github.com/JoseM-18/zincSearch/emailProcessing/processors"
	"github.com/JoseM-18/zincSearch/routes"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	//"runtime/pprof"
	"sync"
)

var emails = make(chan string, 1000)
var dataToZinc = make(chan email.Email, 1000)

// create a struct to store the data

func main() {
	// Create a WaitGroups for the goroutines to wait for each other
	var wgFinders, wgProcessors, wgSender,wg sync.WaitGroup
	//f, _ := os.Create("cpu.pprof")
	//pprof.StartCPUProfile(f)

	wg.Add(1)
	go func() {
		server := os.Getenv("SEARCHING_SERVER_ADDRESS")
		log.Println(http.ListenAndServe(server, routes.SetupRouter()))
	}()

	var rootDirPath string
	flag.StringVar(&rootDirPath, "rootDir", "../allen-p", "path to the root directory")
	flag.Parse()

	res,err := apizinc.CreateIndex()
	if err != nil {
		log.Println(err)
	}
	fmt.Println(res)


	wgFinders.Add(1)
	go finders.FindsDir(rootDirPath, emails, &wgFinders)

	const numProcessors = 4
	// Start the goroutines
	for i := 0; i < numProcessors; i++ {
		wgProcessors.Add(1)
		go processor.ProcessEmails(&wgProcessors, emails, dataToZinc)
	}

	wgSender.Add(1)
	go processor.SendPackages(&wgSender, dataToZinc)

	wgFinders.Wait()
	close(emails)

	wgProcessors.Wait()
	close(dataToZinc)
	wgSender.Wait()

	//defer pprof.StopCPUProfile()
	fmt.Println("Errores Finders: ", finders.GetErroresFinders())
	fmt.Println("Errores Processor: ", processor.GetErroresProcessor())
	fmt.Println("Errores Email: ", email.GetErroresEmail())
	fmt.Println("Errores ApiZinc: ", apizinc.GetErroresApiZinc())
	fmt.Println("Total Errores: ", finders.GetErroresFinders()+processor.GetErroresProcessor()+email.GetErroresEmail()+apizinc.GetErroresApiZinc())

	wg.Wait()
}
