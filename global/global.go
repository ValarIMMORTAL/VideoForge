package global

import "context"

var GlobalCtx, GlobalCancel = context.WithCancel(context.Background())
