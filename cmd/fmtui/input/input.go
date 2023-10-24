package input

import (
	"atomicgo.dev/keyboard"
	"atomicgo.dev/keyboard/keys"
)





func ListenForInput(ch chan keys.Key) error {
	return keyboard.Listen(func(key keys.Key) (stop bool,err error) {
		switch key.Code {
			case keys.RuneKey: {
				switch key.Code.String() {
					case "q": {
						ch<- key
						return true,nil
					}
					default: {
					ch<- key
				return false,nil
		
					}
				}
			}
			case keys.CtrlC,keys.Escape: {
				ch<- key
				return true,nil
			}
			default: {
				ch<- key
				return false,nil
				}
		}
	})
}

