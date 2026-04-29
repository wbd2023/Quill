package styleguide

const (
	eventHeading eventKind = iota + 1
	eventHTMLBlock
	eventListItem
	eventBoundary
)

type eventKind uint8

type documentEvent struct {
	kind     eventKind
	location position
	heading  Heading
	text     string
}

func newHeadingEvent(location position, heading Heading) (event documentEvent) {
	return documentEvent{
		kind:     eventHeading,
		location: location,
		heading:  heading,
	}
}

func newHTMLBlockEvent(location position, text string) (event documentEvent) {
	return documentEvent{
		kind:     eventHTMLBlock,
		location: location,
		text:     text,
	}
}

func newListItemEvent(location position, text string) (event documentEvent) {
	return documentEvent{
		kind:     eventListItem,
		location: location,
		text:     text,
	}
}

func newBoundaryEvent(location position) (event documentEvent) {
	return documentEvent{
		kind:     eventBoundary,
		location: location,
	}
}
