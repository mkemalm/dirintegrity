package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func visit(files *[]string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal(err)
		}
		*files = append(*files, path)
		return nil
	}
}

func addStatsToFile(fileswithstats []string) {
	file, err := os.OpenFile(".stat.out", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}
	defer file.Close()
	var alllines = ""

	for _, line := range fileswithstats {
		alllines += line
	}

	len, err := file.WriteString(alllines)
	if err != nil {
		log.Fatalf("failed writing to file: %s, lenght: %d", err, len)
	}
}

func readStatsFromFile() []string {
	file, err := os.Open(".stat.out")
	if err != nil {
		return nil
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func fileStat(filename string) string {
	out, err := exec.Command("/bin/bash", "-c", "stat -c '%n,%s,%a,%u,%g,%Y,%Z' '"+filename+"'").Output()
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println(string(out))
		fmt.Println("File-->" + filename)
	}
	output := string(out[:])
	return output
}

func removeStatFile() {
	out, err := exec.Command("/bin/bash", "-c", "rm -f .stat.out").Output()
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println(string(out))
	}
}

func trimString(str string) string {
	str = strings.Trim(str, "\r")
	str = strings.Trim(str, " ")
	str = strings.Trim(str, "\n")
	str = strings.Trim(str, "\t")
	return str
}

func checkStats(fstat string, stats []string) bool {

	for _, stat := range stats {
		if trimString(stat) == trimString(fstat) {
			return true
		}
	}

	return false
}

func main() {
	var files []string

	op := os.Args[1]

	root := "/" + os.Args[2]
	err := filepath.Walk(root, visit(&files))
	if err != nil {
		panic(err)
	}

	var lines []string
	for _, file := range files {
		lines = append(lines, fileStat(file))
	}

	if op == "update" {
		removeStatFile()
		addStatsToFile(lines)
	} else if op == "check" {
		stats := readStatsFromFile()
		for _, line := range lines {
			if !checkStats(line, stats) {
				fmt.Println(strings.Split(line, ",")[0] + " is not compliant!")
			}
		}
	}
}
