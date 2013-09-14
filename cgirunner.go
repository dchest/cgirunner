// Copyright 2013 Dmitry Chestnykh. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"log"
	"net/http"
	"net/http/cgi"
	"os"
	"path"
	"path/filepath"
)

var (
	fDir  = flag.String("dir", "", "directory with CGI files")
	fRoot = flag.String("root", "/", "URL prefix")
	fAddr = flag.String("addr", "localhost:8111", "address to run server at")
)

func main() {
	flag.Parse()

	if *fDir == "" {
		flag.Usage()
		os.Exit(1)
	}

	dir, err := os.Open(*fDir)
	if err != nil {
		log.Fatalf("open: %s", err)
	}
	files, err := dir.Readdir(0)
	if err != nil {
		log.Fatalf("readdir: %s", err)
	}
	for _, fi := range files {
		mode := fi.Mode()
		if mode.IsRegular() && (mode.Perm()&0100 != 0) {
			// Register handler.
			urlPath := path.Join(*fRoot, fi.Name())
			filePath := filepath.Join(*fDir, fi.Name())
			http.Handle(urlPath, &cgi.Handler{
				Path: filePath,
				Root: *fRoot,
			})
			log.Printf("Registered %s", urlPath)
		}
	}
	dir.Close()

	log.Printf("serving from %s", *fAddr)
	log.Fatal(http.ListenAndServe(*fAddr, nil))
}
