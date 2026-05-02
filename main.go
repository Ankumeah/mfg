package main

import (
	"errors"
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

const binName = "mfg"
const maxDownloads = 10
const timeout = time.Second * 100
const baseUrl = "https://images.mangafreak.net/downloads/"
const baseMangaUrl = "https://ww2.mangafreak.me/Manga/"

var regex1 *regexp.Regexp
var regex2 *regexp.Regexp
func init() {
  regex1 = regexp.MustCompile(`/Read1_[^"]*`)
  regex2 = regexp.MustCompile(`(?m)[0-9]*$`)
}

var client = &http.Client { Timeout: timeout }

func main() {
  if len(os.Args) < 2 {
    fmt.Printf("Usage: %v <manga>", binName)
    os.Exit(1)
  }

  mangas := os.Args[1:len(os.Args)]

  for _, manga := range mangas {
    fmt.Printf("Downloading %v\n", manga)

    chs, err := getChapters(manga)
    if err != nil {
      fmt.Printf("Error while getting %v: %v\n", manga, err.Error())
      continue
    }

    fmt.Printf("Found %v chapters\n", len(chs))

    var prog uint
    go func() {
      for {
        fmt.Printf("\r%v / %v", len(chs), prog)
        time.Sleep(time.Second * 1)
      }
    }()

    maxPad := 1
    for _, ch := range chs {
      if len(ch) > maxPad { maxPad = len(ch) }
    }

    os.MkdirAll(manga, 0o755)

    var wg sync.WaitGroup
    sem := make(chan struct{}, maxDownloads)
    for _, ch := range chs {
      chInt, err := strconv.Atoi(ch)
      if err != nil {
        fmt.Printf("\nGot non numeric chapter %v: %v\n", ch, err.Error())
        continue
      }

      url := baseUrl + manga + "_" + ch
      savePath := fmt.Sprintf( "%0*d.cbz", maxPad, chInt)

      sem <- struct{}{}
      wg.Go(func() {
        if err := downloadFile(
          url, path.Join(manga, savePath),
        ); err != nil {
          fmt.Printf("\nError while downloading chapter %v: %v\n", ch, err.Error())
        }
        <- sem
        prog++
      })
    }
    wg.Wait()
    fmt.Printf("\r%v / %v", len(chs), prog)
    fmt.Println()
  }
}

func getChapters(manga string) ([]string, error) {
  res, err := client.Get(baseMangaUrl + manga)
  if err != nil { return nil, err }
  defer res.Body.Close()
  html, err := io.ReadAll(res.Body)
  if err != nil { return nil, err }

  match := regex1.FindAllString(string(html), -1)
  match = regex2.FindAllString(strings.Join(match, "\n"), -1)

  seen := make(map[string]struct{})
  var result []string
  for _, ch := range match {
    if _, ok := seen[ch]; ok { continue }

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
  if err != nil { return err }
  defer res.Body.Close()
  if res.StatusCode != http.StatusOK {
    return errors.New("Http error: " + res.Status)
  }

  dir := path.Dir(savePath)
  tmp, err := os.CreateTemp(dir, path.Join(".part_*"))
  if err != nil { return err }
  defer tmp.Close()

  _, err = io.Copy(tmp, res.Body)
  if err != nil { return err }

  err = os.Rename(tmp.Name(), savePath)
  if err != nil { return err }

  return nil
}
