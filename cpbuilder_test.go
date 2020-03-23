/*
 * Copyright (c) 2019.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func setupTestCase(t *testing.T) func(t *testing.T) {
	t.Log("setup test case")

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	testdirpath := filepath.Join(dir, "testprojectdir")
	os.MkdirAll(testdirpath, os.ModePerm)

	dirs := [7]string{"cartridge1", "gradle", "bin", "cartridge2", ".gradle", "config", "cartridge3"}
	for _, d := range dirs {
		os.MkdirAll(filepath.Join(testdirpath, d), os.ModePerm)
	}

	return func(t *testing.T) {
		t.Log("teardown test case")

		os.RemoveAll(testdirpath)
	}
}

func TestParseCommandDirOnly(t *testing.T) {
	os.Args = []string{"command", "--dir=/path/to/dir"}
	m := &Config{}
	m.ParseCommandLine()

	if m.dir != "/path/to/dir" {
		t.Errorf("Directory setting is not ok. It is %s and should be %s", m.dir, "/path/to/dir")
	}
	if m.libPath != "" {
		t.Errorf("Jar path is not correct. It is %s and should be %s", m.libPath, "release/lib ")
	}
	if m.pathExcludes != "" {
		t.Errorf("Excludes configuration is not correct. It is %s and should be %s", m.pathExcludes, "build,target,bin")
	}
}

func TestParseCommand(t *testing.T) {
	os.Args = []string{"command", "--dir=/path/to/dir", "--path=lib/rel", "--excludes=test1,test2"}
	m := &Config{}
	m.ParseCommandLine()

	if m.dir != "/path/to/dir" {
		t.Errorf("Directory setting is not ok. It is %s and should be %s", m.dir, "/path/to/dir")
	}
	if m.libPath != "lib/rel" {
		t.Errorf("Jar path is not correct. It is %s and should be %s", m.pathExcludes, "lib/rel")
	}
	if m.pathExcludes != "test1,test2" {
		t.Errorf("Excludes configuration is not correct. It is %s and should be %s", m.pathExcludes, "test1,test2")
	}
}

func TestDirListing(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	os.Setenv(LIBPATH, "release/lib")
	os.Setenv(PATHEXCLUDES, "build,target,bin,config,gradle")

	m := &Config{}
	m.finalConfigInit()

	m.dir = filepath.Join(dir, "testprojectdir")
	cp := m.createCP()

	if strings.HasPrefix(cp, ":") {
		t.Errorf("Classpath starts with seperator!")
	}

	cpe := strings.Split(cp, string(os.PathListSeparator))
	if len(cpe) != 3 {
		t.Errorf("Classpath contains wrong elements - %d - %s", len(cpe), cp)
	}
}
