/*
 * Copyright 2020 Intershop Communications AG.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	LIBPATH      = "LIBPATH"
	PATHEXCLUDES = "PATHEXCLUDES"
)

var (
	command           = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	paramDir          = command.String("dir", "", "Directory with projects / cartridges")
	paramLibPath      = command.String("path", "", "Sub path with jar file")
	paramPathExcludes = command.String("excludes", "", "String with excluded directories")
)

type Config struct {
	dir          string
	libPath      string
	pathExcludes string
}

func (config *Config) ParseCommandLine() {
	command.Parse(os.Args[1:])

	if *paramDir == "" {
		fmt.Fprintln(os.Stderr, "Parameter 'dir' is empty.")
		command.Usage()
		os.Exit(1)
	}

	config.dir = *paramDir
	config.libPath = *paramLibPath
	config.pathExcludes = *paramPathExcludes
}

func (config *Config) finalConfigInit() {
	envLibPath, availableLibPath := os.LookupEnv(LIBPATH)
	if availableLibPath {
		config.libPath = envLibPath
	}

	envPathExludes, availablePathExludes := os.LookupEnv(PATHEXCLUDES)
	if availablePathExludes {
		config.pathExcludes = envPathExludes
	}
}

func (config *Config) createCP() string {
	files, err := OSReadDir(config.dir)
	if err != nil {
		fmt.Println(err)
	}

	var cp bytes.Buffer
	var jp string
	var excl []string

	if config.libPath != "" {
		jp = config.libPath
	}

	if config.pathExcludes != "" {
		excl = strings.Split(config.pathExcludes, ",")
	}

	for _, dirname := range files {
		if !contains(excl, dirname) {
			if cp.Len() != 0 {
				cp.WriteString(string(os.PathListSeparator))
			}
			cp.WriteString(filepath.Join(config.dir, dirname, jp, "*"))
		}
	}
	return cp.String()
}

func main() {
	config := &Config{}
	config.ParseCommandLine()
	config.finalConfigInit()

	fmt.Println(config.createCP())
}

func OSReadDir(root string) ([]string, error) {
	var files []string
	f, err := os.Open(root)
	if err != nil {
		return files, err
	}
	fileInfo, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return files, err
	}

	for _, file := range fileInfo {
		if file.IsDir() && !strings.HasPrefix(file.Name(), ".") {
			files = append(files, file.Name())
		}
	}
	return files, nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
