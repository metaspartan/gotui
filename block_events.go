package gotui

// HandleEvent handles events. This is a default implementation for Block
// that returns false, meaning the event was not handled.
// Embedders should override this method to handle specific events.
func (b *Block) HandleEvent(e Event) bool {
	return false
}
