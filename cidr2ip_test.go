// Copyright (c) 2023 Roberto Meléndez.
// Licensed under the MIT License. See the LICENSE file in the project root for license information.

package main

import (
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"testing"
)

const binPath = "./cidr2ip"

func TestCmdArguments(t *testing.T) {
	buildBinary(t)

	// Test with a single CIDR as a command-line argument
	file1 := checkCmdOutput(t, binPath, "10.0.0.0/24")

	// Test with multiple CIDRs as command-line arguments
	file2 := checkCmdOutput(t, binPath, "10.0.1.0/24", "172.16.16.0/20", "192.168.0.0/16")

	// Test with no CIDRs provided
	checkError(t, binPath)

	if file1 == file2 {
		removeFiles(t, file1)
	} else {
		removeFiles(t, file1, file2)
	}
}

func TestFileInput(t *testing.T) {
	buildBinary(t)

	// Test with a sample CIDR file from the repository
	file1 := checkCmdOutput(t, binPath, "-f", "cidrs")

	// Test with a non-existent file
	checkError(t, binPath, "-f", "nonexistent_file.txt")

	// Test with an empty file
	file2 := "empty_file.txt"
	if err := createEmptyFile(file2); err != nil {
		t.Fatalf("Failed to create empty file: %v", err)
	}
	checkError(t, binPath, "-f", file2)

	removeFiles(t, file1, file2)
}

func TestInvalidCIDRs(t *testing.T) {
	buildBinary(t)

	// Test with an invalid CIDR as a command-line argument
	checkError(t, binPath, "10.0.0.0/33")

	// Test with multiple CIDRs as command-line arguments, one of which is invalid
	checkError(t, binPath, "10.0.0.0/24", "172.256.0.0/16", "192.168.0.0/16")

	// Test with a file containing one invalid CIDR
	file := "invalid_cidr.txt"
	if err := os.WriteFile(file, []byte("192.168.1.0/16.0"), 0644); err != nil {
		t.Fatalf("Failed to create invalid CIDR file: %v", err)
	}
	checkError(t, binPath, "-f", file)

	removeFiles(t, file)
}

func TestIPCount(t *testing.T) {
	buildBinary(t)

	// Test with a single CIDR as a command-line argument
	file := checkCmdOutput(t, binPath, "172.16.18.0/20")

	data, err := os.ReadFile(file)
	if err != nil {
		t.Fatalf("Error reading file: %v", err)
	}

	count := strings.Count(string(data), "\n")

	expected := 4096
	if count != expected {
		t.Errorf("Expected %d IP addresses, but found %d", expected, count)
	}

	removeFiles(t, file)
}

func buildBinary(t *testing.T) {
	cmd := exec.Command("go", "build", ".")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build binary: %v", err)
	}
}

func checkCmdOutput(t *testing.T, b string, args ...string) string {
	output, err := runCommand(b, args...)
	if err != nil {
		t.Fatalf("Command failed with error: %v", err)
	}

	expected := "IP list saved to cidr2ip_"
	if !strings.Contains(string(output), expected) {
		t.Errorf("Expected '%s', got '%s' instead.", expected, output)
	}

	return extractFileName(output, "IP list saved to (.+\\.csv)")
}

func runCommand(b string, args ...string) (string, error) {
	cmd := exec.Command(b, args...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func extractFileName(output, pattern string) string {
	re := regexp.MustCompile(pattern)
	match := re.FindStringSubmatch(output)
	if len(match) != 2 {
		log.Fatalf("Failed to extract filename from output: %s", output)
	}

	return match[1]
}

func removeFiles(t *testing.T, files ...string) {
	for _, file := range files {
		if err := os.Remove(file); err != nil {
			t.Logf("Error removing file: %v", err)
		}
	}

	if err := os.Remove(binPath); err != nil {
		t.Logf("Failed to remove binary file: %v", err)
	}
}

func createEmptyFile(file string) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()

	return nil
}

func checkError(t *testing.T, b string, args ...string) {
	cmd := exec.Command(b, args...)
	err := cmd.Run()
	if err == nil {
		t.Error("Expected an error, but command succeeded.")
	}
}
