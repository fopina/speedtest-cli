package prober

import "sync"

// Positive number representing some amount of data transfer on a network.
//
type BytesTransferred int64

// Almost an ErrorGroup, but allows individual actions to fail as long as there
// is still some successful transfer. If everything completely fails (no data
// was transferred), the last error received gets returned.
//
type Group struct {
	Grp  sync.WaitGroup
	Sem  chan struct{}
	Inc  chan BytesTransferred
	Errs chan error
	Res  chan BytesTransferred
}

func NewGroup(concurrency int) *Group {
	return &Group{
		Sem:  make(chan struct{}, concurrency),
		Errs: make(chan error),
		Res:  make(chan BytesTransferred),
	}
}

func (p *Group) GetIncremental() chan BytesTransferred {
	if p.Inc == nil {
		p.Inc = make(chan BytesTransferred)
	}
	return p.Inc
}

func (p *Group) Add(probe func() (BytesTransferred, error)) {
	p.Grp.Add(1)
	go func() {
		<-p.Sem
		b, err := probe()
		if err != nil {
			p.Errs <- err
		}
		p.Res <- b
		p.Sem <- struct{}{}
		p.Grp.Done()
	}()
}

func (p *Group) Collect() (BytesTransferred, error) {
	var (
		lastErr   error // Keep the last transfer error in case nothing works.
		totalSize BytesTransferred
		cancel    = make(chan struct{})
	)

	go func() {
		for {
			select {
			case b := <-p.Res:
				totalSize += b
				if p.Inc != nil && totalSize != BytesTransferred(0) {
					p.Inc <- totalSize
				}
			case lastErr = <-p.Errs:
			case _ = <-cancel:
				break
			}
		}

	}()

	for i := 0; i < cap(p.Sem); i++ {
		p.Sem <- struct{}{}
	}
	p.Grp.Wait()
	cancel <- struct{}{}
	for i := 0; i < cap(p.Sem); i++ {
		<-p.Sem
	}

	if p.Inc != nil {
		close(p.Inc)
		p.Inc = nil
	}

	if totalSize != BytesTransferred(0) {
		lastErr = nil
	}
	return totalSize, lastErr
}
