package grok

import (
	"bufio"
	"fmt"
	"os"
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

type RuleSet []Rule

type Rule struct {
	Name    string
	matcher matcher
	skip    bool
}

func (r Rule) Match(pkt []byte, offset int) (Segment, error) {
	bytes, length, err := r.matcher.Match(pkt[offset:])
	if err != nil {
		return Segment{}, err
	}

	return Segment{
		Name:   r.Name,
		Offset: offset,
		Length: length,
		Bytes:  bytes,
		Skip:   r.skip,
	}, nil
}

func LoadFile(filename string) (RuleSet, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	rules := []Rule{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// skip blank lines or comments
		if len(line) == 0 || line[0] == '#' {
			continue
		}

		rule, err := parseRuleLine(line)
		if err != nil {
			return nil, err
		}
		rules = append(rules, rule)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return rules, nil
}

func parseRuleLine(line string) (Rule, error) {
	// split the rule from the name/description
	split := strings.SplitN(line, " ", 2)
	if len(split) != 2 {
		return Rule{}, fmt.Errorf("Invalid description: %q", line)
	}

	ruleStr := split[0]
	description := split[1]

	// grok the rule portion with our gnarly regexp
	match := ruleMatcher.FindStringSubmatch(ruleStr)
	if len(match) == 0 {
		return Rule{}, fmt.Errorf("Invalid rule: %q" + ruleStr)
	}

	// it's a fixed byte rule
	if len(match[2]) > 0 {
		length, _ := strconv.Atoi(match[2])
		return Rule{
			Name:    description,
			matcher: readBytes{length: length},
			skip:    match[1] == "-",
		}, nil
	}

	// it's a scan rule
	targetByte, _ := strconv.ParseUint(match[4], 16, 8)
	return Rule{
		Name:    description,
		matcher: readUntil{targetByte: byte(targetByte)},
		skip:    match[3] == "->",
	}, nil

}
