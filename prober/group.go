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
	grp sync.WaitGroup
	sem chan struct{}
	inc chan BytesTransferred
	err chan error
	res chan BytesTransferred
}

func NewGroup(concurrency int) *Group {
	return &Group{
		sem: make(chan struct{}, concurrency),
		err: make(chan error),
		res: make(chan BytesTransferred),
	}
}

func (p *Group) GetIncremental() chan BytesTransferred {
	if p.inc == nil {
		p.inc = make(chan BytesTransferred)
	}
	return p.inc
}

func (p *Group) Add(probe func() (BytesTransferred, error)) {
	p.grp.Add(1)
	go func() {
		<-p.sem
		b, err := probe()
		if err != nil {
			p.err <- err
		}
		p.res <- b
		p.sem <- struct{}{}
		p.grp.Done()
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
			case b := <-p.res:
				totalSize += b
				if p.inc != nil && totalSize != BytesTransferred(0) {
					p.inc <- totalSize
				}
			case lastErr = <-p.err:
			case <-cancel:
				return
			}
		}

	}()

	for i := 0; i < cap(p.sem); i++ {
		p.sem <- struct{}{}
	}
	p.grp.Wait()
	cancel <- struct{}{}
	for i := 0; i < cap(p.sem); i++ {
		<-p.sem
	}

	if p.inc != nil {
		close(p.inc)
		p.inc = nil
	}

	if totalSize != BytesTransferred(0) {
		lastErr = nil
	}
	return totalSize, lastErr
}
