package debug

type MessageCallback func(topic string, value any)

type MessageHost interface {
	Subscribe(topic string, callback MessageCallback)
}

//##############################################################################

type subscriber struct {
	pattern  string
	callback MessageCallback
}

func (s subscriber) matchTopic(topic string) bool {
	if s.pattern == topic {
		return true
	}
	if s.pattern == "*" {
		return true
	}
	return false
}

type messageHost struct {
	subscribers []subscriber
	messages    map[string]any
}

func (m *messageHost) Subscribe(pattern string, callback MessageCallback) {
	s := subscriber{
		pattern:  pattern,
		callback: callback,
	}

	m.subscribers = append(m.subscribers, s)
	if m.messages != nil {
		for k, v := range m.messages {
			if s.matchTopic(k) {
				s.callback(k, v)
			}
		}
	}
}

func (m *messageHost) Publish(topic string, value any) {
	for _, s := range m.subscribers {
		if s.matchTopic(topic) {
			s.callback(topic, value)
		}
	}
	if m.messages == nil {
		m.messages = make(map[string]any)
	}
	if value == nil {
		delete(m.messages, topic)
	} else {
		m.messages[topic] = value
	}
}
