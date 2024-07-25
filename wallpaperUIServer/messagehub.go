package main

type MessageHub[T any] struct {
	receivers []chan T
}

func CreateMessageHub[T any]() MessageHub[T] {
	return MessageHub[T]{}
}

func (m *MessageHub[T]) WaitForMessage() T {
	c := make(chan T)
	m.receivers = append(m.receivers, c)
	return <-c
}

func (m *MessageHub[T]) WaitForMessageForSelect() chan T {
	c := make(chan T)
	m.receivers = append(m.receivers, c)
	return c
}

func (m *MessageHub[T]) SendMessage(message T) {
	for _, c := range m.receivers {
		c <- message
	}
	m.receivers = []chan T{}
}