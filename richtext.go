package gotgbot

// WalkRichText recursively visits r and every descendant, calling fn on each node.
// Returning false from fn skips that node's children (but not its siblings).
func WalkRichText(r RichText, fn func(RichText) bool) {
	if r == nil {
		return
	}
	if !fn(r) {
		return
	}
	for _, child := range r.Children() {
		WalkRichText(child, fn)
	}
}

// WalkRichBlock recursively visits b and every descendant block, calling fnBlock
// on each. For each block, all RichText fields are also walked via fnText.
// Returning false from fnBlock skips that block's children.
func WalkRichBlock(b RichBlock, fnBlock func(RichBlock) bool, fnText func(RichText) bool) {
	if b == nil {
		return
	}
	if !fnBlock(b) {
		return
	}
	for _, child := range b.RichTextChildren() {
		WalkRichText(child, fnText)
	}
	for _, child := range b.RichBlockChildren() {
		WalkRichBlock(child, fnBlock, fnText)
	}
}

// WalkRichMessage visits every RichBlock and RichText node in a RichMessage.
// Pass nil for either callback to skip that class of node.
func WalkRichMessage(m RichMessage, fnBlock func(RichBlock) bool, fnText func(RichText) bool) {
	if fnBlock == nil {
		fnBlock = func(RichBlock) bool { return true }
	}
	if fnText == nil {
		fnText = func(RichText) bool { return true }
	}
	for _, block := range m.Blocks {
		WalkRichBlock(block, fnBlock, fnText)
	}
}
