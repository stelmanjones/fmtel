package main

import (
	"atomicgo.dev/keyboard"
	"atomicgo.dev/keyboard/keys"
	"github.com/charmbracelet/log"
)

func ListenForInput(ch chan keys.Key) error {
	return keyboard.Listen(func(key keys.Key) (stop bool, err error) {
		switch key.Code {
		case keys.RuneKey:
			{
				switch key.String() {
				case "q":
					return true, nil
				default:

				}
				log.Info("key", key)
				return false, nil
			}
		case keys.CtrlC, keys.Escape:
			{
				ch <- key
				return true, nil
			}
		default:
			{
				ch <- key
				return false, nil
			}
		}
	})
}

func main() {
	ch := make(chan keys.Key)

	ListenForInput(ch)
}
