package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/justmumu/zmapgo"
)

func main() {
	scanner, err := zmapgo.NewAsyncScanner()
	if err != nil {
		log.Fatalf("unable to create async scanner: %v", err)
	}

	err = scanner.AddOptions(
		// Change WithTargets function value
		zmapgo.WithTargets("x.x.x.x/30"),
		zmapgo.WithTargetPort("80"),
		zmapgo.WithRate("10000"),
		zmapgo.WithOutputFile("./output-file.txt"),
		zmapgo.WithLogFile("./log-file.txt"),
	)

	if err != nil {
		log.Fatalf("unable to add options to the async scanner: %v", err)
	}

	if err := scanner.RunAsync(); err != nil {
		log.Fatalf("unable to run async scan: %v", err)
	}

	// Block main until the scan has completed
	if err := scanner.Wait(); err != nil {
		panic(err)
	}

	// it's always good to check for fatals.
	fatals := scanner.GetFatalMessages()
	for _, fatal := range fatals {
		log.Printf("[FATAL]: %s", fatal.Message)
	}
	if len(fatals) > 0 {
		// So zmap did not work as expected and waiting for results would be pointless.
		os.Exit(1)
	}

	infos := scanner.GetInfoMessages()
	for _, info := range infos {
		fmt.Println(strings.Repeat("-", 10))
		fmt.Printf("Log Type: %s\n", info.LogType)
		fmt.Printf("Message: %s\n", info.Message)
	}

	results := scanner.GetResults()
	for _, result := range results {
		fmt.Println(strings.Repeat("-", 10))
		for key, value := range result {
			fmt.Printf("%s: %s\n", key, value)
		}
	}

}
