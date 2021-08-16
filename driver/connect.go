package driver

import (
	"fmt"

	"github.com/driver005/oauth/driver/parse"
	helper "github.com/driver005/oauth/helpers"
	"github.com/gobuffalo/pop/v5"
)

func Init(l *helper.Logger) error {
	pool, idlePool, connMaxLifetime, connMaxIdleTime, cleanedDSN := parse.ParseConnectionOptions(l, "postgres://jntubuvx:D_QoM5kppIE5HjEkhx-bDPOkhftEFMeE@ziggy.db.elephantsql.com/jntubuvx")
	c, err := pop.NewConnection(&pop.ConnectionDetails{
		URL:             parse.FinalizeDSN(l, cleanedDSN),
		IdlePool:        idlePool,
		ConnMaxLifetime: connMaxLifetime,
		ConnMaxIdleTime: connMaxIdleTime,
		Pool:            pool,
	})
	if err != nil {
		return helper.WithStack(err)
	}
	fmt.Println(c)
	return nil
}
