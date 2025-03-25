package global

import "context"

var GlobalCtx, GlobalCancel = context.WithCancel(context.Background())

var OauthCodeChan = make(chan string, 1)
