package go_ws_sh

import (
	"context"
	"errors"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type HertzNext func(c context.Context, r *app.RequestContext)
type HertzMiddleWare func(c context.Context, r *app.RequestContext,next HertzNext)



// Compose 将多个中间件组合成一个中间件
func HertzCompose(middlewares ...HertzMiddleWare) HertzMiddleWare {
    return func(c context.Context, ctx *app.RequestContext, finalNext HertzNext) {
        index := -1
		var dispatch func(i int) HertzNext
        dispatch = func(i int) HertzNext {
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