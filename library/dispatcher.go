package library

type Dispatcher struct {
	Bot      Bot
	updates  chan Update
	handlers *[]Handler
}

func NewDispatcher(bot Bot, updates chan Update) Dispatcher {
	d := Dispatcher{}
	d.Bot = bot
	d.updates = updates
	d.handlers = new([]Handler)
	return d
}

func (d Dispatcher) Start() {
	for upd := range d.updates {
		d.process_update(upd)
	}
}

func (d Dispatcher) process_update(update Update) {
	for _, handler := range *d.handlers {
		if handler.Check_update(update) {
			handler.Handle_update(update, d)
			break
		}
	}
}


func (d Dispatcher) Add_handler(handler Handler) {
	*d.handlers = append(*d.handlers, handler)

}