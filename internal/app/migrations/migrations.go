package migrations

//TODO: create go gen script for migrations

import (
	"fmt"
	"github.com/jackc/pgx"
	"strings"
)

type migrationScript struct {
	before func(*pgx.Conn) error
	sql    string
	after  func(*pgx.Conn) error
}

type migration struct {
	dependsOn string
	version   string
	up        migrationScript
	down      migrationScript
}

type migrations []migration

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
			return fmt.Errorf("all migrations have to have version")
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

func order(m migrations) (migrations, error) {
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

//Migrate applies migrations
func Migrate(pool *pgx.ConnPool) error {
	conn, err := pool.Acquire()
	if err != nil {
		return err
	}

}
