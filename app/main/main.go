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
	"sync"
	//"runtime/pprof"
)

func main() {

	var wg sync.WaitGroup

	wg.Add(1)
	go startHTTPServer()

	var rootDirPath string
	flag.StringVar(&rootDirPath, "rootDir", "../enron_mail_20110402/maildirsd", "path to the root directory")
	flag.Parse()

	initializeAndProcessEmails(rootDirPath)

	printErrorStatistics()

	wg.Wait()

}

/**
 * startHTTPServer starts the HTTP server to listen for requests of searching emails.
 * @returns {void}
 */
func startHTTPServer() {
		server := os.Getenv("SEARCHING_SERVER_ADDRESS")
	log.Println(http.ListenAndServe(server, routes.SetupRouter()))
}

/**
 * initializeAndProcessEmails initializes the search engine and processes the emails.
 * @param {string} rootDirPath - The path to the root directory.
 * @returns {void}
 */
func initializeAndProcessEmails(rootDirPath string) {
	//f, _ := os.Create("cpuProfile2.pprof")
			//pprof.StartCPUProfile(f)

	
	// create a WaitGroups for the goroutines to wait for each other
	var wgFinders, wgProcessors, wgSender sync.WaitGroup
	
	// Create channels to send the data between the goroutines
	var dirsEmails = make(chan string, 1000)
	var dataToZinc = make(chan email.Email, 1000)
	
	// create index
	res, err := apizinc.CreateIndex()
	if err != nil {
		log.Println(err)
	}
	fmt.Println(res)
	
	// start finders to find emails files and send them to the emails channel
	wgFinders.Add(1)
	go finders.FindsDir(rootDirPath, dirsEmails, &wgFinders)

	//map to store uniques dates
	dateMap := &sync.Map{}

	
	// start processors to process emails and send them to the dataToZinc channel
	const numProcessors = 4
	for i := 0; i < numProcessors; i++ {
		wgProcessors.Add(1)
		go processor.ProcessEmails(&wgProcessors, dirsEmails, dataToZinc, dateMap)
	}

	// start sender to send the data to the search engine
	wgSender.Add(1)
	go processor.SendPackages(&wgSender, dataToZinc)

	// Wait for the goroutines to finish
	wgFinders.Wait()
	close(dirsEmails)
	wgProcessors.Wait()
	close(dataToZinc)
	wgSender.Wait()

	//pprof.StopCPUProfile()
	//f.Close()
}

/**
 * printErrorStatistics prints the number of errors that occurred when searching for email messages,
 * parsing them, and sending them to the search engine.
 * @returns {void}
 */
func printErrorStatistics() {
	fmt.Println("Errores Finders: ", finders.GetErroresFinders())
	fmt.Println("Errores Processor: ", processor.GetErroresProcessor())
	fmt.Println("Errores Email: ", email.GetErroresEmail())
	fmt.Println("Errores ApiZinc: ", apizinc.GetErroresApiZinc())
	fmt.Println("Total Errores: ", finders.GetErroresFinders()+processor.GetErroresProcessor()+email.GetErroresEmail()+apizinc.GetErroresApiZinc())
}
