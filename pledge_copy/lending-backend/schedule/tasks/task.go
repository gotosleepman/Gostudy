package tasks

import (
	"lending-copy/db"
	"lending-copy/schedule/services"
	"time"

	"github.com/jasonlvhit/gocron"
)

func Task() {
	err := db.RedisFlushDB()
	if err != nil {
		panic("clear redis error " + err.Error())
	}
	services.NewPool().UpdateAllPoolInfo()
	services.NewBalanceMonitor().Monitor()

	s := gocron.NewScheduler()
	s.ChangeLoc(time.UTC)
	_ = s.Every(2).Minutes().From(gocron.NextTick()).Do(services.NewPool().UpdateAllPoolInfo)
	_ = s.Every(30).Minutes().From(gocron.NextTick()).Do(services.NewBalanceMonitor().Monitor)
	<-s.Start()
}
