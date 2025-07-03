package wechat

import (
	"github.com/WeChat-Easy-Chat/route"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

func init() {
	functions.HTTP("wechat", route.URL)
}
