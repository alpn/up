package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/ttacon/chalk"
	"github.com/vbauerster/mpb"
	"github.com/vbauerster/mpb/decor"
)

const total = 100

func updateBar(bar *mpb.Bar, name string) {

	stat := B2Client.Status().Writers[name]
	if stat != nil {
		prog := stat.Progress
		var sum float64
		for _, p := range prog {
			sum += p
		}
		bar.SetTotal(int64(total*len(prog)), false)
		bar.SetCurrent(int64(total*sum) + 1)
	}
}

func showProgress(stop chan bool, wg *sync.WaitGroup, bucketName string, name string) {

	writerName := fmt.Sprintf("%s/%s", bucketName, name)
	p := mpb.New(mpb.WithWidth(64))

	bar := p.AddBar(int64(total),
		mpb.PrependDecorators(

			// display name with one space on the right
			decor.Name(name, decor.WC{W: len(name) + 1, C: decor.DidentRight}),

			// replace ETA decorator with "done" message, OnComplete event
			decor.OnComplete(
				decor.AverageETA(decor.ET_STYLE_GO, decor.WC{W: 4}), fmt.Sprintf("%sdone%s", chalk.Green, chalk.Reset),
			),
		),
		mpb.AppendDecorators(decor.Percentage()),
	)

	defer func() {
		bar.SetTotal(1, true)
		bar.SetCurrent(1)
		p.Wait()
		wg.Done()
	}()

	for {

		updateBar(bar, writerName)

		timeoutchan := make(chan bool)
		go func() {
			<-time.After(time.Millisecond * 25)
			timeoutchan <- true
		}()

		select {
		case <-timeoutchan:
		case <-stop:
			return
		}
	}
}
