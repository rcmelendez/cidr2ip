// Copyright (c) 2023 Roberto Meléndez.
// Licensed under the MIT License. See the LICENSE file in the project root for license information.

package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"net"
	"os"
	"time"
)

const (
	app     = "cidr2ip"
	version = "1.0.0"
)

func main() {
	var (
		fileFlag    string
		helpFlag    bool
		versionFlag bool
	)

	flag.StringVar(&fileFlag, "f", "", "Specify a `filename` with CIDRs")
	flag.BoolVar(&helpFlag, "h", false, "Show help menu")
	flag.BoolVar(&versionFlag, "v", false, "Show version")
	flag.Parse()

	if versionFlag {
		printVersion()
		os.Exit(0)
	}

	if helpFlag {
		printHelp()
		os.Exit(0)
	}

	if fileFlag == "" && flag.NArg() == 0 {
		fmt.Fprintln(os.Stderr, "Error: No CIDRs provided. Use -h for help.")
		os.Exit(1)
	}

	cidrs, err := readCIDRs(fileFlag)
	handleError(err)

	ips, err := generateIPs(cidrs)
	handleError(err)

	file := fmt.Sprintf("%s_%s.csv", app, time.Now().Format("2006-01-02_15-04-05"))
	err = saveToCSV(ips, file)
	handleError(err)

	fmt.Printf("IP addresses saved to %s", file)
}

func printHelp() {
	fmt.Printf("Usage: %s [-f filename] <CIDR1 CIDR2 ...>\nOptions:\n", app)
	flag.PrintDefaults()
}

func printVersion() {
	fmt.Printf("%s version %s\n", app, version)
}

func readCIDRs(file string) ([]string, error) {
	if file != "" {
		return readFromFile(file)
	}

	return flag.Args(), nil
}

func readFromFile(file string) ([]string, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cidrs []string

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		cidrs = append(cidrs, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return cidrs, nil
}

func generateIPs(cidrs []string) ([]string, error) {
	ips := []string{}
	for _, cidr := range cidrs {
		ipList, err := getIPsFromCIDR(cidr)
		if err != nil {
			return nil, err
		}
		ips = append(ips, ipList...)
	}

	return ips, nil
}

func getIPsFromCIDR(cidr string) ([]string, error) {
	ips := []string{}

	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); nextIP(ip) {
		ips = append(ips, ip.String())
	}

	return ips, nil
}

func nextIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func saveToCSV(ips []string, file string) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	for _, ip := range ips {
		err := w.Write([]string{ip})
		if err != nil {
			return err
		}
	}

	return nil
}

func handleError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
