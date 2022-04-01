package migrations

import "github.com/jackc/pgx"

var initialMigration = migration{
	version:   "1",
	dependsOn: "",
	up: func(conn *pgx.Conn) error {
		sql :=
			`
				CREATE TABLE IF NOT EXISTS "companies"
				(
					uuid VARCHAR(36)  not null UNIQUE,
					name VARCHAR(100) not null,
					PRIMARY KEY (uuid)
				);
				
				CREATE TABLE IF NOT EXISTS "projects"
				(
					uuid         VARCHAR(36)  not null UNIQUE,
					name         VARCHAR(100) not null,
					company_uuid VARCHAR(36)  not null,
					PRIMARY KEY (uuid),
					CONSTRAINT fk_projects_companies
						FOREIGN KEY (company_uuid)
							REFERENCES companies (uuid) ON DELETE CASCADE ON UPDATE CASCADE
				);
				
				CREATE TYPE token_type AS ENUM ('bearer', 'jwt');
				
				CREATE TABLE IF NOT EXISTS "projects_settings"
				(
					project_uuid     VARCHAR(36) not null UNIQUE,
					token_expiration SERIAL,
					token_type       token_type  not null DEFAULT 'jwt',
					logout_enabled   BOOL        not null DEFAULT TRUE,
					PRIMARY KEY (project_uuid),
					CONSTRAINT fk_projects_settings_projects
						FOREIGN KEY (project_uuid)
							REFERENCES projects (uuid) ON DELETE CASCADE ON UPDATE CASCADE
				);
				
				CREATE TABLE IF NOT EXISTS "users"
				(
					uuid         VARCHAR(36)  not null UNIQUE,
					username     VARCHAR(100) not null UNIQUE,
					password     VARCHAR(500) not null,
					project_uuid VARCHAR(36)  not null,
					PRIMARY KEY (uuid),
					CONSTRAINT fk_users_projects
						FOREIGN KEY (project_uuid)
							REFERENCES projects (uuid) ON DELETE CASCADE ON UPDATE CASCADE
				);
				
				CREATE TABLE IF NOT EXISTS "roles"
				(
					name VARCHAR(50) not null UNIQUE,
					PRIMARY KEY (name)
				);
			`
		_, err := conn.Exec(sql)
		return err
	},
	down: func(conn *pgx.Conn) error {
		sql := `
		DROP TABLE IF EXISTS 'users';
		DROP TABLE IF EXISTS 'projects_settings';
		DROP TABLE IF EXISTS 'projects';
		DROP TABLE IF EXISTS 'companies';
		DROP TABLE IF EXISTS 'roles';
		`
		_, err := conn.Exec(sql)
		return err
	},
}
