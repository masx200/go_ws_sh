package go_ws_sh

import (
	"context"
	"errors"
	//"log"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"github.com/masx200/go_ws_sh/types"
)

//仿照koa-compose的写法
// 定义一个中间件类型，它接受一个context和一个*app.RequestContext作为参数
// 并返回一个函数，该函数也接受一个context和一个*app.RequestContext作为参数

// Compose 将多个中间件组合成一个中间件
func HertzCompose(middlewares ...types.HertzMiddleWare) types.HertzMiddleWare {
	return func(c context.Context, ctx *app.RequestContext, finalNext types.HertzNext) {
		// defer func() {
		// 	if r := recover(); r != nil {
		// 		log.Println(r)
		// 		ctx.AbortWithMsg("Internal Server Error", consts.StatusInternalServerError)
		// 	}
		// }()
		if len(middlewares) == 1 {
			middlewares[0](c, ctx, finalNext)
			return
		}

		if len(middlewares) == 0 {
			finalNext(c, ctx)
			return
		}
		index := -1
		var dispatch func(i int) types.HertzNext
		dispatch = func(i int) types.HertzNext {
			return func(c context.Context, ctx *app.RequestContext) {
				if i <= index {
					ctx.AbortWithError(
						consts.StatusInternalServerError,

						errors.New("next() called multiple times"))
					return
				}
				index = i
				if i == len(middlewares) {
					finalNext(c, ctx)
					return
				}
				middlewares[i](c, ctx, dispatch(i+1))
			}
		}
		dispatch(0)(c, ctx)
	}
}
