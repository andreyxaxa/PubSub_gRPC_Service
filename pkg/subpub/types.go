package subpub

import (
	"context"
	"fmt"
	"sync"

	subpuberr "github.com/andreyxaxa/PubSub_gRPC_Service/pkg/subpub/errors"
)

type subscriber struct {
	handler   MessageHandler   // вызываем при получении сообщения
	ch        chan interface{} // сюда приходят сообщения
	closeCh   chan struct{}    // канал для остановки горутин
	stopOnce  sync.Once        // остановка только один раз
	unsubFunc func()           // отписка
}

func (s *subscriber) start() {
	// в отдельной горутине
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Subscriber - start() - panic recovered")
			}
		}()

		for {
			select {
			// читаем из канала
			case msg, ok := <-s.ch:
				if !ok {
					return
				}
				// если прочли - отдаем хендлеру
				s.handler(msg)
			case <-s.closeCh:
				return
			}
		}
	}()
}

func (s *subscriber) stop() {
	// хватит читать
	s.stopOnce.Do(
		func() {
			close(s.closeCh)
		})
}

// subPubImpl

type subPubImpl struct {
	mu          sync.RWMutex
	subscribers map[string]map[*subscriber]struct{} // [тема]множество_подписчиков
	closed      bool
	closeCh     chan struct{}
}

// Создает нового подписчика.
// Запускает горутину для получения и обработки сообщений.
func (s *subPubImpl) Subscribe(subject string, cb MessageHandler) (Subscription, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return nil, subpuberr.ErrSubPubClosed
	}

	sub := &subscriber{
		ch:      make(chan interface{}, 16),
		handler: cb,
		closeCh: make(chan struct{}),
	}

	sub.start()

	if s.subscribers[subject] == nil {
		s.subscribers[subject] = make(map[*subscriber]struct{})
	}
	s.subscribers[subject][sub] = struct{}{}

	sub.unsubFunc = func() {
		s.mu.Lock()
		defer s.mu.Unlock()
		delete(s.subscribers[subject], sub)
		sub.stop()
	}

	// Unsubscribe()
	return &subscription{stop: sub.unsubFunc}, nil
}

// Находит всех подписчиков по теме.
// Отправляет сообщение в их канал.
func (s *subPubImpl) Publish(subject string, msg interface{}) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.closed {
		return subpuberr.ErrSubPubClosed
	}

	for sub := range s.subscribers[subject] {
		select {
		case sub.ch <- msg:
			// delivered
		case <-sub.closeCh:
			// subscriber closed
		}
	}

	return nil
}

// Помечает систему как закрытую.
// Останавливает всех подписчиков.
// Ожидаем завершения всех горутин, пока не истечет контекст.
func (s *subPubImpl) Close(ctx context.Context) error {
	s.mu.Lock()

	if s.closed {
		s.mu.Unlock()
		return nil
	}

	s.closed = true
	close(s.closeCh)

	wg := &sync.WaitGroup{}

	for _, subs := range s.subscribers {
		for sub := range subs {
			wg.Add(1)
			go func(sub *subscriber) {
				defer wg.Done()
				sub.stop()
			}(sub)
		}
	}
	s.mu.Unlock()

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
		return nil
	}
}

// subscription

type subscription struct {
	stop func()
}

func (s *subscription) Unsubscribe() {
	s.stop()
}
