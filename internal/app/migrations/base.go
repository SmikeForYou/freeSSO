package migrations

//TODO: create go gen script for migrations

import (
	"errors"
	"fmt"
	"freeSSO/internal/app/connections"
	"freeSSO/internal/app/logger"
	"strings"

	"github.com/jackc/pgx"
	"github.com/sirupsen/logrus"
)

type migration struct {
	dependsOn string
	version   string
	up        func(*pgx.Conn) error
	down      func(*pgx.Conn) error
}

type migrations []migration

var log *logrus.Entry = logger.GetNamedLogger("migrations")

func (m migrations) getMigration(version string) (migration, error) {
	for _, mg := range m {
		if mg.version == version {
			return mg, nil
		}
	}
	return migration{}, fmt.Errorf("migration %s not found", version)
}

func (m migrations) getMigrationByDependency(depVersion string) (migration, error) {
	for _, mg := range m {
		if mg.dependsOn == depVersion {
			return mg, nil
		}
	}
	return migration{}, fmt.Errorf("migration with dependency '%s' not found", depVersion)
}

func validateVersions(m migrations) error {
	for _, mg := range m {
		if mg.version == "" {
			return fmt.Errorf("all migrations need version")
		}
	}
	return nil
}

func validateDependencies(m migrations) error {
	independentVersions := make([]string, 0)
	for _, mg := range m {
		if mg.dependsOn == "" {
			independentVersions = append(independentVersions, mg.version)
		}
	}
	if len(independentVersions) > 2 {
		return fmt.Errorf("only one(initial) migration can be independent. "+
			"Please check dependsOn value in migrations: %s", strings.Join(independentVersions, ", "))
	}
	if len(independentVersions) == 0 {
		return fmt.Errorf("seems that initial migration(migration without dependency) does not exists")
	}
	return nil
}

func (m migrations) order() (migrations, error) {
	nmg := make(migrations, 0)
	for i := range m {
		if i == 0 {
			initMigration, err := m.getMigrationByDependency("")
			if err != nil {
				return nil, err
			}
			nmg = append(nmg, initMigration)
			continue
		}
		prevMg := nmg[i-1]
		mg, err := m.getMigrationByDependency(prevMg.version)
		if err != nil {
			return nil, err
		}
		nmg = append(nmg, mg)
	}
	return nmg, nil
}

func (m migrations) getHead(version string) migrations {
	res := make(migrations, 0)
	for _, m := range m {
		res = append(res, m)
		if m.version == version {
			return res
		}
	}
	return m
}
func (m migrations) apply(version string, pool *pgx.ConnPool) error {
	mig := m.getHead(version)
	conn, err := pool.Acquire()
	if err != nil {
		return err
	}
	for _, mi := range mig {
		err = mi.up(conn)
		if err != nil {
			return err
		}
	}
	return nil

}
func (m migrations) getLatest() string {
	ordered, _ := m.order()
	return ordered[len(ordered)-1].version
}

//Migrate applies migrations
func MigrateWithPool(migr migrations, pool *pgx.ConnPool) error {
	err := validateVersions(migr)
	if err != nil {
		log.Error(errors.New("error validating migrations versions"))
		log.Error(err)
	}
	err = validateDependencies(migr)
	if err != nil {
		log.Error(errors.New("error validating migrations dependencies"))
		log.Error(err)
	}
	ordered, err := migr.order()
	if err != nil {
		log.Error(errors.New("error ordering migrations"))
		log.Error(err)
	}
	err = ordered.apply(ordered.getLatest(), pool)
	return err
}

func Migrate() error {
	pool := connections.GetDbConnPool()
	err := MigrateWithPool(migrationsList, pool)
	if err != nil {
		log.Error(err)
	}
	return err
}
