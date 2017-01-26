package transport

// writer provides resource calls a writer to respond to the proto request. This
// type impliments io.WriteCloser.
type writer struct {
	resv  chan []byte
	close chan bool
}

// Writes b to the writer
func (t *writer) Write(b []byte) (int, error) {
	t.resv <- b
	return len(b), nil
}

// Close closes the writer
func (t *writer) Close() error {
	select {
	case <-t.close:
		return nil
	default:
		close(t.close)
		return nil
	}
}
