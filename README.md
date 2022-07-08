<p align="center">
    <img width="400" src="images/logo.png">
<p>

<p align="center">
	<a href="LICENSE">
		<img src="https://img.shields.io/badge/License-MIT-yellow.svg">
	</a>
	<a href="https://pkg.go.dev/github.com/justmumu/zmapgo">
		<img src="https://pkg.go.dev/badge/github.com/justmumu/zmapgo.svg" alt="Go Reference">
	</a>
	<a href="https://goreportcard.com/report/github.com/justmumu/zmapgo">
        <img src="https://goreportcard.com/badge/github.com/justmumu/zmapgo">
    </a>
	<a href="https://github.com/justmumu/zmapgo/actions/workflows/test.yml">
		<img src="https://github.com/justmumu/zmapgo/actions/workflows/test.yml/badge.svg">
	</a>
	<a href="https://github.com/justmumu/zmapgo/actions/workflows/build.yml">
		<img src="https://github.com/justmumu/zmapgo/actions/workflows/build.yml/badge.svg">
	</a>
	<a href='https://coveralls.io/github/justmumu/zmapgo?branch=main'>
		<img src='https://coveralls.io/repos/github/justmumu/zmapgo/badge.svg?branch=main' alt='Coverage Status' />
	</a>
<p>

This library aims to provide to golang developers an idiomatic interface for zmap version 2.1.1.

Inspired by the [nmap](https://github.com/Ullaakut/nmap) library.

## What is Zmap
[Zmap](https://github.com/zmap/zmap) is a network tool for scanning the entire Internet (or large samples). ZMap is capable of scanning the entire Internet in around 45 minutes on a gigabit network connection, reaching ~98% theoretical line speed.

[More Details](https://github.com/zmap/zmap/wiki)

## Supported Features
- [x] All of `zmap 2.1.1` native options.
- [x] Cancellable contexts support
- [x] Validation for options
- [x] Async Scanner
- [x] Blocking Scanner

## TODO
- [ ] More examples

## Installation
```bash
go get github.com/justmumu/zmapgo
```

## Simple Example
```go
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
    // Create Context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

    // Create Blocking Scanner with InitOptions
    scanner, err := zmapgo.NewBlockingScanner(
		zmapgo.WithContext(ctx),
	)
	if err != nil {
		log.Fatalf("unable to create zmap scanner: %v", err)
	}
    
    // Add Options to scanner
    // Equivalent to `zmap ---target-port 80 1.1.1.0/30 --rate 10000 --output-fields saddr,sport --log-file ./log-file.txt --output-file ./output-file.txt`
    err = scanner.AddOptions(
		zmapgo.WithTargets("1.1.1.0/30"),
		zmapgo.WithTargetPort("80"),
		zmapgo.WithRate("10000"),
        zmapgo.WithOutputFields([]string{"saddr", "sport"}),
		zmapgo.WithLogFile("./log-file.txt"),
		zmapgo.WithOutputFile("./output-file.txt"),
	)
	if err != nil {
		log.Fatalf("unable to add options: %v", err)
	}

    // Run the scan
    results, _, _, _, _, fatals, err := scanner.RunBlocking()
	if err != nil {
		log.Fatalf("unable to run zmap scan: %v", err)
	}

    // It's always good to check for fatals.
	if len(fatals) > 0 {
		// So zmap did not work as expected and waiting for results would be pointless.
		for _, fatal := range fatals {
			log.Printf("[FATAL]: %s", fatal.Message)
		}
		os.Exit(1)
	}

    // Print All Results
	for _, result := range results {
		fmt.Printf("%s\n", strings.Repeat("-", 20))
		for key, value := range result {
			fmt.Printf("%s: %s\n", key, value)
		}
	}
```

The program output:

```bash
--------------------
saddr: 1.1.1.3
sport: 80
--------------------
saddr: 1.1.1.1
sport: 80
--------------------
saddr: 1.1.1.0
sport: 80
--------------------
saddr: 1.1.1.2
sport: 80
```
## LICENCE
This project is under [MIT License](LICENSE)