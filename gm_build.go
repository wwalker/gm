package main

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// buildMd compiles the infile (xxx.md | stdin) to outfile (xxx.html | stdout)
func buildMd(infile string) {
	// get the dir for link replacement, if any
	dir := filepath.Dir(infile)
	// Get the input
	var input io.Reader
	if infile != "" {
		f, err := os.Open(infile)
		if err != nil {
			check(err, "Problem opening", infile)
			return
		}
		defer f.Close()
		input = f
	} else {
		input = os.Stdin
		dir = "."
	}
	// read the input
	markdown, err := ioutil.ReadAll(input)
	check(err, "Problem reading the markdown.")

	//compile the input
	html, err := compile(markdown)
	check(err, "Problem compiling the markdown.")
	if localmdlinks {
		html = replaceLinks(html, dir)
	}

	// output the result
	if infile == "" {
		os.Stdout.Write(html)
	} else {
		outfile := filepath.Join(outdir, infile[:len(infile)-3]+".html")
		if os.MkdirAll(filepath.Dir(outfile), os.ModePerm) != nil {
			check(err, "Problem to reach/create folder:", filepath.Dir(outfile))
		}
		err = ioutil.WriteFile(outfile, html, 0644)
		check(err, "Problem modifying", outfile)
	}
}

// buildFiles convert all .md files verifying one of the patterns to .html
func buildFiles() {
	// check all patterns
	for _, pattern := range inpatterns {
		info("Looking for '%s'.\n", pattern)
		// if the input is piped
		if pattern == "stdin" {
			buildMd("")
			continue
		}
		// look for all files with the given patterns
		// but build only .md ones
		allfiles, err := filepath.Glob(pattern)
		check(err, "Problem looking for file pattern:", pattern)
		if len(allfiles) == 0 {
			info("No files found.\n")
			continue
		}
		for _, infile := range allfiles {
			if strings.HasSuffix(infile, ".md") {
				info("  Converting %s...", infile)
				buildMd(infile)
				info("done.\n")
			}
		}
	}
}
