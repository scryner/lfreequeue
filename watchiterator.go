package lfreequeue

type watchIterator struct {
	queue *Queue
	quit chan int
}

func (w *watchIterator) Iter() <-chan interface{} {
	c := make(chan interface{})
	go w.watchAndIterate(c)
	return c
}

func (w *watchIterator) watchAndIterate(c chan<- interface{}) {
	for {
		datum, ok := queue.Dequeue()

		if !ok {
			notify := queue.WatchWakeup()

			select {
			case <-notify:
				continue
			case <-quit:
				goto endIteration
			}

			<-queue.WatchWakeup()
		} else {
			c <- datum
		}
	}

endIteration:
	close(c)
}

func (w *watchIterator) Close() {
	w.quit <- 1
}