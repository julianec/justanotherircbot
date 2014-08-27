package main

import (
	"container/ring"
	"github.com/julianec/go-ircevent"
	"time"
)

type bufferedMessage struct {
	Channel string
	Message string
}

type MessageBuffer struct {
	buffer   *ring.Ring
	ircconn  *irc.Connection
	maxsize  int32
	lostmsgs int
	ticker   *time.Ticker
	karma    int
}

func NewMessageBuffer(ircconn *irc.Connection, size int32) *MessageBuffer {
	var buff *MessageBuffer = &MessageBuffer{
		buffer:   ring.New(1),
		ircconn:  ircconn,
		maxsize:  size,
		lostmsgs: 0,
		ticker:   time.NewTicker(2 * time.Second),
		karma:    0,
	}
	go buff.deliverMessages()
	go buff.refillKarma()
	return buff
}

func (m *MessageBuffer) AddMessage(channel, message string) {
	if int32(m.buffer.Len()) >= m.maxsize { // buffer full?
		m.lostmsgs += 1
		return
	}
	var elem *ring.Ring = ring.New(1)
	var msg *bufferedMessage = &bufferedMessage{
		Channel: channel,
		Message: message,
	}
	elem.Value = msg
	// Add element elem to buffer.
	m.buffer.Prev().Link(elem)
}

func (m *MessageBuffer) deliverMessages() {
	for { // forever
		if m.karma == 0 { // karma good enough?
			for m.karma > -5 && m.buffer.Next() != m.buffer { // karma left  AND not empty
				var cur *ring.Ring = m.buffer.Unlink(1)
				var ok bool
				var message *bufferedMessage

				message, ok = cur.Value.(*bufferedMessage)
				if ok {
					// Write message from buffer to channel and update karma
					m.ircconn.Privmsg(message.Channel, message.Message)
					m.karma -= 1
				}
			}
		}
		// no messages left or empty karma, sleep for 2 seconds before trying it again
		time.Sleep(2 * time.Second)
	}
}

func (m *MessageBuffer) refillKarma() {
	for _ = range m.ticker.C { //increment karma every 2 seconds
		if m.karma < 0 {
			m.karma += 1
		}
	}
}
