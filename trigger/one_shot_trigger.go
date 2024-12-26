package trigger

type OneShotTrigger struct{}

func (t *OneShotTrigger) Listen() chan bool {
	nextChan := make(chan bool)

	go func() {
		nextChan <- true
		close(nextChan)
	}()

	return nextChan
}
