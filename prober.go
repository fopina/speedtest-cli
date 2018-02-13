package speedtest

import (
	"sync"
	"time"
)

type bytesTransferred int

type proberGroup struct {
	concurrency int
	wg          sync.WaitGroup

	sem         chan struct{}
	incremental chan BytesPerSecond
	errors      chan error
	results     chan bytesTransferred
}

func newProberGroup(concurrency int) *proberGroup {
	pg := &proberGroup{
		concurrency: concurrency,

		sem:     make(chan struct{}, concurrency),
		errors:  make(chan error),
		results: make(chan bytesTransferred),
	}

	// Load up the semaphore so things block initially.
	for i := 0; i < concurrency; i++ {
		pg.sem <- struct{}{}
	}

	return pg
}

func (p *proberGroup) GetIncrementalResults() chan BytesPerSecond {
	if p.incremental == nil {
		p.incremental = make(chan BytesPerSecond)
	}
	return p.incremental
}

func (p *proberGroup) Add(probe func() (bytesTransferred, error)) {
	p.wg.Add(1)
	go func() {
		p.sem <- struct{}{}

		b, err := probe()
		if err != nil {
			p.errors <- err
		}
		p.results <- b

		p.wg.Done()
		<-p.sem
	}()
}

func (p *proberGroup) Collect() (BytesPerSecond, error) {
	go func() {
		p.wg.Wait()
		close(p.results)
		close(p.errors)
	}()

	var (
		lastErr   error // Keep the last transfer error in case nothing works.
		totalSize bytesTransferred
		start     = time.Now()
	)

	calc := func() BytesPerSecond {
		return BytesPerSecond(float64(totalSize) * float64(time.Second) /
			float64(time.Since(start)))
	}

	// Fill the semaphore, starting workgroup tasks.
	for i := 0; i < p.concurrency; i++ {
		<-p.sem
	}

	go func() {
		for err := range p.errors {
			lastErr = err
		}
	}()

	for b := range p.results {
		totalSize += b
		if p.incremental != nil && totalSize != bytesTransferred(0) {
			p.incremental <- calc()
		}
	}

	if totalSize == bytesTransferred(0) {
		return BytesPerSecond(0), lastErr
	} else {
		return calc(), nil
	}
}
