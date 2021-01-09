// Command telegrup is the command line telegra.ph file uploader.
// Run with `-h` to get some help.
package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rusq/telegraph"
)

var (
	list    = flag.String("l", "", "read the list of files from the text `file`.\nEach line should contain one file.")
	skip    = flag.Bool("s", false, "skip failed uploads")
	quiet   = flag.Bool("q", false, "be quiet (errors are printed anyway)")
	timeout = flag.Duration("t", 60*time.Second, "single file upload `timeout`")
	js      = flag.Bool("json", false, "results as json")
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Upload files to telegra.ph.\n\nUsage: %s [flags] < -l <file> | <filename>... >\n", filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}
}

type result struct {
	n    int
	path string
	err  error
}

func main() {
	flag.Parse()

	files, err := getFileList(flag.Args(), *list)
	if err != nil {
		flag.Usage()
		log.Fatal(err)
	}

	results, err := uploadBunch(files, *skip, *timeout)
	if err != nil {
		log.Fatal(err)
	}
	if *quiet {
		return
	} else {
		// output results
		if err := printResults(os.Stdout, results, *js); err != nil {
			log.Fatal(err)
		}
	}
}

func getFileList(files []string, listfile string) ([]string, error) {
	if len(files) == 0 && listfile == "" {
		return nil, errors.New("no files provided")
	}
	if len(files) != 0 {
		return files, nil
	}
	f, err := os.Open(listfile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		filename := strings.TrimSpace(scanner.Text())
		if filename == "" {
			continue
		}
		files = append(files, filename)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("no data discovered in %s", listfile)
	}
	return files, nil
}

// uploadBunch uploads a bunch of files, returning results.
func uploadBunch(files []string, skip bool, timeout time.Duration) ([]result, error) {
	var results = make([]result, 0, len(files))

	for i, filename := range files {
		remotePath, err := uploadOne(filename, timeout)
		if err != nil {
			msg := fmt.Sprintf("error uploading file %d: %s : %s", i+1, filename, err)
			if !skip {
				return nil, errors.New(msg) // OUCH
			}
			log.Print("SKIPPED: " + msg)
		}
		results = append(results, result{n: i, path: remotePath, err: err})
	}
	return results, nil
}

// usually telegra.ph shouldnt return more than one upload result given one
// file. BUT YOU NEVER KNOW.
const usuallyReturned = 1

// uploadOne uploads just one file.
func uploadOne(filename string, timeout time.Duration) (string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	res, err := telegraph.Upload(ctx, f)
	if err != nil {
		return "", err
	}
	if n := len(res); n != usuallyReturned {
		return "", fmt.Errorf("unexpected number of results: %d", n)
	}
	return res[0].Src, err
}

// printResults prints the results to writer.
func printResults(w io.Writer, results []result, asJson bool) error {
	if asJson {
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		return enc.Encode(&results)
	}
	for _, res := range results {
		if res.err != nil {
			if _, err := fmt.Fprintf(w, "%2d: ERROR: %s", res.n, res.err); err != nil {
				return err
			}
			continue
		}
		if _, err := fmt.Fprintf(w, "%2d: OK: %s%s\n", res.n, telegraph.BaseURL, res.path); err != nil {
			return err
		}
	}
	return nil
}
