package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("domain,hasMX,hasSPF,spfRecord,hasDWARC,dwarcRecord")

	for scanner.Scan() {
		checkDomain(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("ERROR: could not read from input: %v\n", err)
	}
}

func checkDomain(domain string) {
	var hasMX, hasDWARC, hasSPF bool
	var spfRecord, dwarcRecord string

	mxRecords, err := net.LookupMX(domain)
	if err != nil {
		log.Printf("Error looking up MX records for %s: %v\n", domain, err)
	}

	if len(mxRecords) > 0 {
		hasMX = true
	}

	txtRecords, err := net.LookupTXT(domain)
	if err != nil {
		log.Printf("Error looking up TXT records for %s: %v\n", domain, err)
	}

	for _, record := range txtRecords {
		if strings.HasPrefix(record, "v=spf1") {
			hasSPF = true
			spfRecord = record
			break
		}
	}

	dmarcRecords, err := net.LookupTXT("_dmarc." + domain)
	if err != nil {
		log.Printf("Error looking up DMARC records for %s: %v\n", domain, err)
	}

	for _, dmarc := range dmarcRecords {
		if strings.HasPrefix(dmarc, "v=DMARC1") {
			hasDWARC = true
			dwarcRecord = dmarc
			break
		}
	}

	fmt.Printf("%v,%v,%v,%v,%v,%v\n", domain, hasMX, hasSPF, spfRecord, hasDWARC, dwarcRecord)

	if hasSPF && hasMX && hasDWARC {
		fmt.Println("Valid Site :)")
	} else {
		fmt.Println("Not a Valid Site :(")
	}
}
