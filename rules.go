package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

/*
| Rule | Description |
| ---- | ----------- |
| N    | Read N bytes |
| -N   | Read N bytes, but don't display them |
| >XX  | Read until seeing XX |
| ->XX | Read until seeing XX, but don't display it |
*/

var ruleMatcher *regexp.Regexp

type rule struct {
	Description string
	Matcher     matcher
	Skip        bool
}

type matcher interface {
	Match([]byte) (matched []byte, length int, err error)
}

func init() {
	/*
		Group 1: Negation for N bytes
		Group 2: N bytes
		Group 3: Read/Skip marker for scan
		Group 4: Target byte for scan
	*/
	ruleMatcher = regexp.MustCompile("(-)?([0-9]+)|(->|>)([0-9a-fA-F]{1,2})")
}

func loadRules(filename string) ([]rule, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	rules := []rule{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// skip blank lines or comments
		if len(line) == 0 || line[0] == '#' {
			continue
		}

		split := strings.SplitN(line, " ", 2)
		if len(split) != 2 {
			return nil, fmt.Errorf("Invalid description: %q", line)
		}

		ruleStr := split[0]
		description := split[1]

		match := ruleMatcher.FindStringSubmatch(ruleStr)
		log.Printf("%q", ruleStr)
		for i, m := range match {
			log.Printf("%d : %q", i, m)
		}
		if len(match) == 0 {
			return nil, fmt.Errorf("Invalid rule: %q" + ruleStr)
		}

		// it's a fixed byte rule
		if len(match[2]) > 0 {
			length, _ := strconv.Atoi(ruleStr)
			rules = append(rules, rule{
				Description: description,
				Matcher:     readBytes{length: length},
				Skip:        match[1] == "-",
			})
			continue
		}

		// it's a scan rule
		targetByte, _ := strconv.ParseUint(match[4], 16, 8)
		rules = append(rules, rule{
			Description: description,
			Matcher:     scanToByte{targetByte: byte(targetByte)},
			Skip:        match[3] == "->",
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return rules, nil
}

type readBytes struct {
	length int
}

func (r readBytes) Match(pkt []byte) ([]byte, int, error) {
	bytes := make([]byte, r.length)
	copy(bytes, pkt[:r.length])
	return bytes, r.length, nil
}

type scanToByte struct {
	targetByte byte
}

func (r scanToByte) Match(pkt []byte) ([]byte, int, error) {
	indexOf := bytes.IndexByte(pkt, r.targetByte)
	if indexOf == -1 {
		return nil, 0, fmt.Errorf("Could not find terminating byte %02x", r.targetByte)
	}

	bytes := make([]byte, indexOf+1)
	copy(bytes, pkt[:indexOf+1])
	return bytes, indexOf + 1, nil
}
