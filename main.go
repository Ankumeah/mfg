package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

const baseUrl = "https://images.mangafreak.net/downloads/"
const baseMangaUrl = "https://ww2.mangafreak.me/Manga/"

var maxDownloads uint
var timeout time.Duration
var mangas []string

var regex1 *regexp.Regexp
var regex2 *regexp.Regexp

func init() {
	if len(os.Args) < 1 {
		fmt.Fprintln(os.Stderr, "argv is of len 0. Your runnning this on a potato arent you?")
	}

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %v [OPTIONS]... [MANGAS]...\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr, "  --help -h")
		fmt.Fprintln(os.Stderr, "\tPrint this text")
	}

	flag.UintVar(&maxDownloads, "j", 10, "Max parallel downloads")
	t := flag.String("t", "10s", "Max download timeout. Eg - 1h1m1s1ms1us1ns")
	flag.Parse()
	mangas = flag.Args()

	if len(mangas) < 1 {
		flag.Usage()
		fmt.Fprintln(os.Stderr, "Must provide atlest one manga")
		os.Exit(1)
	}

	var err error
	timeout, err = time.ParseDuration(*t)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid timeout: %v\n", err.Error())
		os.Exit(1)
	}

	regex1 = regexp.MustCompile(`/Read1_[^"]*`)
	regex2 = regexp.MustCompile(`(?m)[0-9]*$`)
}

var client = &http.Client{Timeout: timeout}

func main() {
	for _, manga := range mangas {
		fmt.Printf("Downloading %v\n", manga)

		chs, err := getChapters(manga)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error while getting %v: %v\n", manga, err.Error())
			continue
		}

		maxPad := 1
		for i, ch := range chs {
			if ch == "" {
				chs = append(chs[:i], chs[i+1:]...)
			}
			if len(ch) > maxPad {
				maxPad = len(ch)
			}
		}

		fmt.Printf("Found %v chapters\n", len(chs))

		var prog uint
		go func() {
			for {
				fmt.Fprintf(os.Stderr, "\r%v / %v", len(chs), prog)
				time.Sleep(time.Second * 1)
			}
		}()

		os.MkdirAll(manga, 0o755)

		var wg sync.WaitGroup
		sem := make(chan struct{}, maxDownloads)
		for _, ch := range chs {
			chInt, err := strconv.Atoi(ch)
			if err != nil {
				fmt.Fprintf(os.Stderr, "\nGot non numeric chapter %v: %v\n", ch, err.Error())
				continue
			}

			url := baseUrl + manga + "_" + ch
			savePath := fmt.Sprintf("%0*d.cbz", maxPad, chInt)

			sem <- struct{}{}
			wg.Go(func() {
				if err := downloadFile(
					url, path.Join(manga, savePath),
				); err != nil {
					fmt.Fprintf(os.Stderr, "\nError while downloading chapter %v: %v\n", ch, err.Error())
				}
				<-sem
				prog++
			})
		}
		wg.Wait()
		fmt.Fprintf(os.Stderr, "\r%v / %v", len(chs), prog)
		fmt.Fprintln(os.Stderr)
	}
}

func getChapters(manga string) ([]string, error) {
	res, err := client.Get(baseMangaUrl + manga)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	html, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	match := regex1.FindAllString(string(html), -1)
	match = regex2.FindAllString(strings.Join(match, "\n"), -1)

	seen := make(map[string]struct{})
	var result []string
	for _, ch := range match {
		if _, ok := seen[ch]; ok {
			continue
		}

		seen[ch] = struct{}{}
		result = append(result, ch)
	}

	return result, nil
}

func downloadFile(url string, savePath string) error {
	if _, err := os.Stat(savePath); !errors.Is(err, os.ErrNotExist) {
		return nil
	} else if !errors.Is(err, os.ErrNotExist) && err != nil {
		return err
	}

	res, err := client.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return errors.New("Http error: " + res.Status)
	}

	tmp, err := os.CreateTemp("", "mfg_*")
	if err != nil {
		return err
	}
	defer tmp.Close()

	_, err = io.Copy(tmp, res.Body)
	if err != nil {
		return err
	}

	err = os.Rename(tmp.Name(), savePath)
	if err != nil {
		return err
	}

	return nil
}
