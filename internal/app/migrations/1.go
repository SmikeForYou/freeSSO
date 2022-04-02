package migrations

import (
	"context"

	"github.com/jackc/pgx/v4"
)

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
				
				CREATE TYPE internal_permissions AS ENUM ('project', 'company', 'all_companies');
				CREATE TABLE IF NOT EXISTS "users"
				(
					uuid         VARCHAR(36)  not null UNIQUE,
					email     	 VARCHAR(100) not null UNIQUE,
					password     VARCHAR(500) not null,
					project_uuid VARCHAR(36)  not null,
					internal_permissions 	internal_permissions default null;
					PRIMARY KEY (uuid),
					CONSTRAINT fk_users_projects
						FOREIGN KEY (project_uuid)
							REFERENCES projects (uuid) ON DELETE CASCADE ON UPDATE CASCADE
				);
			`
		_, err := conn.Exec(context.TODO(), sql)
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
		_, err := conn.Exec(context.TODO(), sql)
		return err
	},
}

var fieldsMigration = migration{
	dependsOn: "1",
	version:   "2",
	up: func(conn *pgx.Conn) error {
		sql := `
			CREATE TABLE IF NOT EXISTS "fields"
			(
				type VARCHAR(60) NOT NULL UNIQUE,
				PRIMARY KEY (type)
			);
			CREATE TABLE IF NOT EXISTS "validators" (
				name VARCHAR(60) NOT NULL UNIQUE,
				PRIMARY KEY (name)
			);
			CREATE TABLE IF NOT EXISTS "field_validators" (
				field_type VARCHAR(60) NOT NULL,
				validator VARCHAR(60) NOT NULL,
				PRIMARY KEY (field_type, validator)
			);
			CREATE TABLE IF NOT EXISTS "user_fields" (
				uuid VARCHAR(36) NOT NULL,
				user_uuid VARCHAR(36) NOT NULL,
				field_name VARCHAR(60) NOT NULL,
				field_type VARCHAR(60) NOT NULL,
				PRIMARY KEY (uuid),
				CONSTRAINT fk_users_fields
					FOREIGN KEY (user_uuid)
						REFERENCES users (uuid) ON DELETE CASCADE ON UPDATE CASCADE,
				CONSTRAINT fk_fields_field_types
					FOREIGN KEY (field_type)
						REFERENCES fields (type) ON DELETE CASCADE ON UPDATE CASCADE,
			);
			CREATE TABLE IF NOT EXISTS "user_fields_validators" (
				user_field_uuid VARCHAR(36) NOT NULL,
				field_name VARCHAR(60) NOT NULL,
				PRIMARY KEY (user_field_uuid), 
				
			)
		`
		_, err := conn.Exec(context.TODO(), sql)
		return err
	},
	down: func(conn *pgx.Conn) error {
		sql := `
			DROP TABLE IF EXISTS "fields";
			DROP TABLE IF EXISTS "validators";
			DROP TABLE IF EXISTS "user_fields";
		`
		_, err := conn.Exec(context.TODO(), sql)
		return err
	},
}
