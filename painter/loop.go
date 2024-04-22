package painter

import (
	"image"
	"sync"

	"golang.org/x/exp/shiny/screen"
)

// Receiver отримує текстуру, яка була підготовлена в результаті виконання команд у циклі подій.
type Receiver interface {
	Update(t screen.Texture)
}

// Loop реалізує цикл подій для формування текстури отриманої через виконання операцій отриманих з внутрішньої черги.
type Loop struct {
	Receiver Receiver

	buffer screen.Texture // текстура, яка зараз формується

	mq messageQueue

	stop    chan struct{}
	stopReq bool
}

var size = image.Pt(800, 800)

// Start запускає цикл подій. Цей метод потрібно запустити до того, як викликати на ньому будь-які інші методи.
func (l *Loop) Start(s screen.Screen) {
	l.buffer, _ = s.NewTexture(size)
	l.stop = make(chan struct{})

	go func() {
		for !(l.stopReq && l.mq.empty()) {

			op := l.mq.pull()

			update := op.Do(l.buffer)

			if update {
				l.Receiver.Update(l.buffer)
			}
		}
		close(l.stop)
	}()
}

// Post додає нову операцію у внутрішню чергу.
func (l *Loop) Post(op Operation) {
	l.mq.push(op)
}

// StopAndWait сигналізує про необхідність завершити цикл та блокується до моменту його повної зупинки.
func (l *Loop) StopAndWait() {
	l.Post(OperationFunc(func(screen.Texture) {
		l.stopReq = true
	}))

	<-l.stop
}

// Черга повідомлень, яка використовується для зберігання операцій.
type messageQueue struct {
	ops    []Operation
	mu     sync.Mutex
	signal chan struct{}
}

func (mq *messageQueue) push(op Operation) {
	mq.mu.Lock()
	defer mq.mu.Unlock()

	mq.ops = append(mq.ops, op)

	if mq.signal != nil {
		close(mq.signal)
	}
}

func (mq *messageQueue) pull() Operation {
	mq.mu.Lock()
	defer mq.mu.Unlock()

	if len(mq.ops) == 0 {
		mq.mu.Unlock()

		mq.signal = make(chan struct{})
		<-mq.signal
		mq.signal = nil

		mq.mu.Lock()
	}
	op := mq.ops[0]
	mq.ops[0] = nil
	mq.ops = mq.ops[1:]
	return op
}

func (mq *messageQueue) empty() bool {
	mq.mu.Lock()
	defer mq.mu.Unlock()

	return len(mq.ops) == 0
}
