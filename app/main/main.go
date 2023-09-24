package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime/pprof"
	"sync"
	apizinc "github.com/JoseM-18/zincSearch/apiZinc"
	"github.com/JoseM-18/zincSearch/email"
	"github.com/JoseM-18/zincSearch/emailProcessing/finders"
	processor "github.com/JoseM-18/zincSearch/emailProcessing/processors"
	"github.com/JoseM-18/zincSearch/routes"
)

var emails = make(chan string, 1000)
var dataToZinc = make(chan email.Email, 1000)

// create a struct to store the data

func main() {
	// Create a WaitGroups for the goroutines to wait for each other
	var wgFinders, wgProcessors, wgSender sync.WaitGroup
	f, _ := os.Create("cpu.pprof")
	pprof.StartCPUProfile(f)

	go func() {
		server := os.Getenv("SEARCHING_SERVER_ADDRESS")
		log.Println(http.ListenAndServe(server, routes.SetupRouter()))
	}()

	var rootDirPath string
	flag.StringVar(&rootDirPath, "rootDir", "../enron_mail_20110402/maildirsda", "path to the root directory")
	flag.Parse()

	apizinc.CreateIndex()

	wgFinders.Add(1)
	go finders.FindsDir(rootDirPath, emails, &wgFinders)

	// Start the goroutines
	for i := 0; i < 4; i++ {
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

	defer pprof.StopCPUProfile()
	fmt.Println("Errores Finders: ", finders.GetErroresFinders())
	fmt.Println("Errores Processor: ", processor.GetErroresProcessor())
	fmt.Println("Errores Email: ", email.GetErroresEmail())
	fmt.Println("Errores ApiZinc: ", apizinc.GetErroresApiZinc())
	fmt.Println("Total Errores: ", finders.GetErroresFinders()+processor.GetErroresProcessor()+email.GetErroresEmail()+apizinc.GetErroresApiZinc())


}
