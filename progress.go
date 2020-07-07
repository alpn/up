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

func makePipeBar(name string) (*mpb.Progress, *mpb.Bar) {
	p := mpb.New(mpb.WithWidth(14))
	bar := p.AddSpinner(int64(total), mpb.SpinnerOnMiddle,
		mpb.SpinnerStyle([]string{"∙∙∙∙", "●∙∙∙", "∙●∙∙", "∙∙●∙", "∙∙∙●", "∙∙∙∙"}),
		mpb.PrependDecorators(
			decor.Name(name),
		),
		mpb.AppendDecorators(
			decor.OnComplete(
				decor.EwmaETA(decor.ET_STYLE_GO, 60),
				fmt.Sprintf(" %sdone%s ", chalk.Green, chalk.Reset),
			),
		),
	)
	return p, bar
}

func makeFileBar(name string) (*mpb.Progress, *mpb.Bar) {
	p := mpb.New(mpb.WithWidth(64))
	bar := p.AddBar(int64(total),
		mpb.PrependDecorators(
			decor.Name(name, decor.WC{W: 1, C: decor.DidentRight}),
			decor.OnComplete(
				decor.AverageETA(decor.ET_STYLE_GO, decor.WC{W: 4}),
				fmt.Sprintf(" %sdone%s ", chalk.Green, chalk.Reset),
			),
		),
		mpb.AppendDecorators(
			decor.Name(chalk.Magenta.String(), decor.WC{W: 0}),
			decor.Percentage(),
			decor.Name(chalk.Reset.String(), decor.WC{W: 0})),
	)
	return p, bar
}
func showProgress(stop chan bool, wg *sync.WaitGroup, bucketName string, name string, isPipe bool) {

	writerName := fmt.Sprintf("%s/%s", bucketName, name)

	var p *mpb.Progress
	var bar *mpb.Bar

	if isPipe {
		p, bar = makePipeBar(name)
	} else {
		p, bar = makeFileBar(name)
	}

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
