package gormx

import (
	"github.com/prometheus/client_golang/prometheus"
	"gorm.io/gorm"
	"time"
)

type Callbacks struct {
	vector *prometheus.SummaryVec
}

func (c *Callbacks) Name() string {
	return "prometheus_callback"
}

func NewCallbacks(opts prometheus.SummaryOpts) *Callbacks {
	vector := prometheus.NewSummaryVec(opts, []string{"type", "table"})
	prometheus.MustRegister(vector)
	return &Callbacks{vector: vector}
}

func (c *Callbacks) Initialize(db *gorm.DB) error {
	// 监控创建
	err := db.Callback().Create().Before("*").Register("prometheus_create_before", c.before())

	if err != nil {
		panic(err)
	}
	err = db.Callback().Create().After("*").Register("prometheus_create_after", c.after("create"))
	if err != nil {
		panic(err)
	}

	// 监控更新
	err = db.Callback().Update().Before("*").Register("prometheus_update_before", c.before())

	if err != nil {
		panic(err)
	}
	err = db.Callback().Update().After("*").Register("prometheus_update_after", c.after("update"))
	if err != nil {
		panic(err)
	}

	// 监控删除
	err = db.Callback().Delete().Before("*").Register("prometheus_delete_before", c.before())

	if err != nil {
		panic(err)
	}
	err = db.Callback().Delete().After("*").Register("prometheus_delete_after", c.after("delete"))
	if err != nil {
		panic(err)
	}
	// 监控查询
	err = db.Callback().Query().Before("*").Register("prometheus_query_before", c.before())
	if err != nil {
		panic(err)
	}
	err = db.Callback().Query().After("*").Register("prometheus_query_after", c.after("query"))
	if err != nil {
		panic(err)
	}

	err = db.Callback().Row().Before("*").Register("prometheus_row_before", c.before())
	if err != nil {
		panic(err)
	}
	err = db.Callback().Row().After("*").Register("prometheus_row_after", c.after("row"))
	if err != nil {
		panic(err)
	}

	err = db.Callback().Raw().Before("*").Register("prometheus_raw_before", c.before())
	if err != nil {
		panic(err)
	}
	err = db.Callback().Raw().After("*").Register("prometheus_raw_after", c.after("raw"))
	if err != nil {
		panic(err)
	}
	return nil

}
func (c *Callbacks) before() func(db *gorm.DB) {
	return func(db *gorm.DB) {
		startTime := time.Now()
		db.Set("start_time", startTime)
	}
}

func (c *Callbacks) after(typ string) func(db *gorm.DB) {
	return func(db *gorm.DB) {
		val, _ := db.Get("start_time")
		startTime := val.(time.Time)
		table := db.Statement.Table
		if table == "" {
			table = "unknown"
		}

		c.vector.WithLabelValues("create", table).Observe(float64(time.Since(startTime).Milliseconds()))
		//l.Info("执行时间", logx.Any("cost", time.Since(startTime).Milliseconds()))
	}
}
