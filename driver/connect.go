package driver

import (
	"fmt"

	"github.com/gobuffalo/pop/v5"
	"github.com/ory/x/sqlcon"

	"github.com/ory/x/logrusx"

	"github.com/ory/x/errorsx"
)

func Init(l *logrusx.Logger) error {
	pool, idlePool, connMaxLifetime, connMaxIdleTime, cleanedDSN := sqlcon.ParseConnectionOptions(l, "postgres://jntubuvx:D_QoM5kppIE5HjEkhx-bDPOkhftEFMeE@ziggy.db.elephantsql.com/jntubuvx")
	c, err := pop.NewConnection(&pop.ConnectionDetails{
		URL:             sqlcon.FinalizeDSN(l, cleanedDSN),
		IdlePool:        idlePool,
		ConnMaxLifetime: connMaxLifetime,
		ConnMaxIdleTime: connMaxIdleTime,
		Pool:            pool,
	})
	if err != nil {
		return errorsx.WithStack(err)
	}
	fmt.Println(c)
	return nil
}
