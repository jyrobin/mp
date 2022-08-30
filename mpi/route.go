package mpi

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jyrobin/goutil"
	"github.com/jyrobin/mp"
)

func Mount(g *gin.RouterGroup, dom mp.Domain) { // , logging jog.Logging) {
	if goutil.IsNil(dom) {
		panic("Nil domain")
	}

	g.POST("/call", func(c *gin.Context) {
		var req Request
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		mthd, m, opts := req.Unpack()
		ctx := c.Request.Context()
		//logging.For(ctx).With(
		//	zap.String("method", mthd),
		//	zap.String("kind", m.Kind()),
		//	zap.String("ns", m.Ns()),
		//).Info("Call")

		fmt.Println(mthd, m.Json("  "))

		if ret, err := dom.Call(ctx, mthd, m, opts...); err != nil {
			//logging.For(ctx).With(
			//	zap.String("method", mthd),
			//	zap.String("kind", m.Kind()),
			//	zap.String("ns", m.Ns()),
			//).Error(err.Error())

			fmt.Println(err)

			me, code := Error(err, 400)
			c.JSON(code, me)
		} else {
			//logging.For(ctx).With(
			//	zap.String("kind", ret.Kind()),
			//	zap.String("ns", ret.Ns()),
			//	zap.Int("list", len(ret.List())),
			//).Info("Result")

			c.JSON(200, ret)
		}
	})

	/* Later
	g.POST("/list", func(c *gin.Context) {
		var req Request
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		mthd, m, opts := req.Unpack()
		ctx := c.Request.Context()
		logging.For(ctx).With(
			zap.String("method", mthd),
			zap.String("kind", m.Kind()),
			zap.String("ns", m.Ns()),
		).Info("List meta")

		if ret, err := dom.List(ctx, m, opts...); err != nil {
			c.JSON(400, goutil.SimpleJsonError(err.Error(), 400))
		} else {
			logging.For(ctx).With(
				zap.Int("size", len(ret.List())),
			).Info("Got meta list")
			c.JSON(200, ret)
		}
	})

	g.POST("/first", func(c *gin.Context) {
		var req Request
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		m, opts := req.ToMeta()
		ctx := c.Request.Context()
		logging.For(ctx).With(
			zap.String("kind", m.Kind()),
			zap.String("ns", m.Ns()),
		).Info("First meta")

		if ret, err := dom.First(ctx, m, opts...); err != nil {
			c.JSON(400, goutil.SimpleJsonError(err.Error(), 400))
		} else {
			logging.For(ctx).With(
				zap.String("kind", ret.Kind()),
				zap.String("ns", ret.Ns()),
				zap.String("gid", ret.Gid()),
			).Info("Got meta")
			c.JSON(200, ret)
		}
	})

	g.POST("/actor", func(c *gin.Context) {
		var req Request
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		m, opts := req.ToMeta()
		ctx := c.Request.Context()
		logging.For(ctx).With(
			zap.String("kind", m.Kind()),
			zap.String("ns", m.Ns()),
		).Info("First actor")

		if actor, err := dom.FirstActor(ctx, m, opts...); err != nil {
			c.JSON(400, goutil.SimpleJsonError(err.Error(), 400))
		} else {
			am := actor.Meta()
			logging.For(ctx).With(
				zap.String("meta", am.Label()),
			).Info("Got actor")
			c.JSON(200, am)
		}
	})

	g.POST("/process", func(c *gin.Context) {
		var req Request
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		m, opts := req.ToMeta()
		ctx := c.Request.Context()
		logging.For(ctx).With(
			zap.String("kind", m.Kind()),
			zap.String("ns", m.Ns()),
			zap.String("gid", m.Gid()),
		).Info("Do actor")

		if actor, err := dom.FirstActor(ctx, m, opts...); err != nil {
			c.JSON(400, goutil.SimpleJsonError(err.Error(), 400))
		} else {
			logging.For(ctx).With(
				zap.String("kind", actor.Meta().Kind()),
				zap.String("ns", actor.Meta().Ns()),
				zap.String("gid", actor.Meta().Gid()),
			).Info("Do actor")

			ret, err := actor.Process(ctx, m)
			if err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
			} else {
				c.JSON(200, ret)
			}
		}
	})

	g.POST("/do", func(c *gin.Context) {
		var req Request
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		m, opts := req.ToMeta()
		ctx := c.Request.Context()
		logging.For(ctx).With(
			zap.String("kind", m.Kind()),
			zap.String("ns", m.Ns()),
			zap.String("gid", m.Gid()),
		).Info("Do actor")

		if actor, err := dom.FirstActor(ctx, m, opts...); err != nil {
			c.JSON(400, goutil.SimpleJsonError(err.Error(), 400))
		} else {
			logging.For(ctx).With(
				zap.String("kind", actor.Meta().Kind()),
				zap.String("ns", actor.Meta().Ns()),
				zap.String("gid", actor.Meta().Gid()),
			).Info("Do actor")

			ret, err := actor.Process(ctx, m)
			if err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
			} else {
				c.JSON(200, ret)
			}
		}
	})
	*/
}
