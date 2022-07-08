package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/justmumu/zmapgo"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	scanner, err := zmapgo.NewBlockingScanner(
		zmapgo.WithContext(ctx),
	)
	if err != nil {
		log.Fatalf("unable to create zmap scanner: %v", err)
	}

	err = scanner.AddOptions(
		// Change WithTargets function value
		zmapgo.WithTargets("x.x.x.x/30"),
		zmapgo.WithTargetPort("80"),
		zmapgo.WithRate("10000"),
		zmapgo.WithLogFile("./log-file.txt"),
		zmapgo.WithOutputFile("./output-file.txt"),
	)
	if err != nil {
		log.Fatalf("unable to add options: %v", err)
	}

	results, _, _, _, infos, fatals, err := scanner.RunBlocking()
	if err != nil {
		log.Fatalf("unable to run zmap scan: %v", err)
	}

	// it's always good to check for fatals.
	if len(fatals) > 0 {
		// So zmap did not work as expected and waiting for results would be pointless.
		for _, fatal := range fatals {
			log.Printf("[FATAL]: %s", fatal.Message)
		}
		os.Exit(1)
	}

	// Print Info Messages.
	// It is also an example for all trace, debug, warning, info and fatal messages.
	for _, info := range infos {
		fmt.Printf("%s\n", strings.Repeat("-", 20))
		fmt.Printf("Log Type: %s\n", info.LogType)
		fmt.Printf("Log Time: %s\n", info.LogTime.Format("01-02 15:04:05"))
		fmt.Printf("Messge: %s\n", info.Message)
	}

	// Print All Results
	for _, result := range results {
		fmt.Printf("%s\n", strings.Repeat("-", 20))
		for key, value := range result {
			fmt.Printf("%s: %s\n", key, value)
		}
	}
}
