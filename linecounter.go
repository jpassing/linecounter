//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//

package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s directory glob...\n", os.Args[0])
		os.Exit(1)
	}

	directory := os.Args[1]
	globs := os.Args[2:]
	totalLines := 0

	filepath.Walk(
		directory,
		func(path string, info os.FileInfo, err error) error {
			for _, glob := range globs {
				match, err := filepath.Match(glob, info.Name())
				if err == nil && match {
					lines, err := CountLines(path)
					if err != nil {
						return err
					}

					fmt.Printf("%-80s: %d\n", path, lines)

					totalLines += lines
				}
			}

			return nil
		})

	fmt.Printf("\nTotal                                                                           : %d\n", totalLines)
}

func CountLines(path string) (int, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, err
	}

	scanner := bufio.NewScanner(f)
	lines := 0
	for scanner.Scan() {
		// Count non-whitespace, non-comment lines
		sourceLine := strings.TrimSpace(scanner.Text())
		if len(sourceLine) > 0 && !strings.HasPrefix(sourceLine, "//") {
			lines++
		}
	}

	f.Close()

	return lines, nil
}
