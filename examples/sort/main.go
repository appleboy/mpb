package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/appleboy/mpb"
	"github.com/appleboy/mpb/decor"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	var wg sync.WaitGroup
	p := mpb.New(mpb.WithWaitGroup(&wg))
	total := 100
	numBars := 3
	wg.Add(numBars)

	for i := 0; i < numBars; i++ {
		var name string
		if i != 1 {
			name = fmt.Sprintf("Bar#%d:", i)
		}
		b := p.AddBar(int64(total),
			mpb.PrependDecorators(
				decor.StaticName(name, 0, decor.DwidthSync),
				decor.CountersNoUnit("%d / %d", 10, decor.DSyncSpace),
			),
			mpb.AppendDecorators(
				decor.ETA(3, 0),
			),
		)
		go func() {
			defer wg.Done()
			max := 200 * time.Millisecond
			for i := 0; i < total; i++ {
				time.Sleep(time.Duration(rand.Intn(10)+1) * max / 10)
				if i&1 == 1 {
					priority := total - int(b.Current())
					p.UpdateBarPriority(b, priority)
				}
				b.Increment()
			}
		}()
	}

	p.Wait()
	fmt.Println("done")
}
