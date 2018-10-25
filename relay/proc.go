package relay

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/codec"
)

// 处理入站的relay消息
func ResoleveInboundEvent(inputEvent cellnet.Event) (ouputEvent cellnet.Event, handled bool, err error) {

	switch relayMsg := inputEvent.Message().(type) {
	case *RelayACK:

		ev := &RecvMsgEvent{
			Ses: inputEvent.Session(),
			ack: relayMsg,
		}

		if relayMsg.MsgID != 0 {

			ev.Msg, _, err = codec.DecodeMessage(int(relayMsg.MsgID), relayMsg.Msg)
			if err != nil {
				return
			}
		}

		if log.IsDebugEnabled() {

			peerInfo := inputEvent.Session().Peer().(cellnet.PeerProperty)

			log.Debugf("#relay.recv(%s)@%d len: %d %s passThrough: '%+v' | %s",
				peerInfo.Name(),
				inputEvent.Session().ID(),
				cellnet.MessageSize(ev.Message()),
				cellnet.MessageToName(ev.Message()),
				relayMsg.PassThrough(),
				cellnet.MessageToString(ev.Message()))
		}

		if bcFunc != nil {
			// 转到对应线程中调用
			cellnet.SessionQueuedCall(inputEvent.Session(), func() {
				bcFunc(ev)
			})
		}

		return ev, true, nil
	}

	return inputEvent, false, nil
}

// 处理relay.Relay出站消息的日志
func ResolveOutboundEvent(inputEvent cellnet.Event) (handled bool, err error) {

	switch relayMsg := inputEvent.Message().(type) {
	case *RelayACK:
		if log.IsDebugEnabled() {

			var payload interface{}
			if relayMsg.MsgID != 0 {

				payload, _, err = codec.DecodeMessage(int(relayMsg.MsgID), relayMsg.Msg)
				if err != nil {
					return
				}
			}

			peerInfo := inputEvent.Session().Peer().(cellnet.PeerProperty)

			log.Debugf("#relay.send(%s)@%d len: %d %s passThrough: '%+v' | %s",
				peerInfo.Name(),
				inputEvent.Session().ID(),
				cellnet.MessageSize(payload),
				cellnet.MessageToName(payload),
				relayMsg.PassThrough(),
				cellnet.MessageToString(payload))

			return true, nil

		}

	}

	return
}
