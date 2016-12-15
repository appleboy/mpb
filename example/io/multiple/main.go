package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/vbauerster/mpb"
)

func main() {
	log.SetOutput(os.Stderr)

	url1 := "https://homebrew.bintray.com/bottles/youtube-dl-2016.12.12.sierra.bottle.tar.gz"
	url2 := "https://homebrew.bintray.com/bottles/libtiff-4.0.7.sierra.bottle.tar.gz"

	p := mpb.New().SetWidth(60)

	for i, url := range [...]string{url1, url2} {
		p.Wg.Add(1) // if you omit this line, main will return without waiting for download goroutines
		name := fmt.Sprintf("url%d:", i+1)
		go download(p, name, url)
	}

	p.WaitAndStop()
	fmt.Println("Finished")
}

func download(p *mpb.Progress, name, url string) {
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("%s: %v", name, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("non-200 status: %s", resp.Status)
		log.Printf("%s: %v", name, err)
		return
	}

	size := resp.ContentLength

	// create dest
	destName := filepath.Base(url)
	dest, err := os.Create(destName)
	if err != nil {
		err = fmt.Errorf("Can't create %s: %v", destName, err)
		log.Printf("%s: %v", name, err)
		return
	}

	// create bar with appropriate decorators
	bar := p.AddBar(int(size)).
		PrependCounters(mpb.UnitBytes, 20).
		PrependName(name, len(name)).
		AppendETA()
	// create proxy reader
	reader := bar.ProxyReader(resp.Body)
	// and copy from reader
	_, err = io.Copy(dest, reader)

	if closeErr := dest.Close(); err == nil {
		err = closeErr
	}
	if err != nil {
		log.Printf("%s: %v", name, err)
	}
}