package repository

import (
	"database/sql"
	"github.com/GuiaBolso/darwin"
)

func Migrate(db *sql.DB) error {
	driver := darwin.NewGenericDriver(db, darwin.PostgresDialect{})
	d := darwin.New(driver, migrations, nil)
	return d.Migrate()
}

var migrations = []darwin.Migration{
	{
		Version:     1,
		Description: "Initial",
		Script: `
			create extension if not exists "uuid-ossp";

			create table if not exists users
			(
				id         uuid                     default uuid_generate_v4() not null
					constraint users_pkey
						primary key,
				first_name varchar(30)                                         not null,
				last_name  varchar(150)                                        not null,
				email      varchar(255)                                        not null,
				password   varchar(255)                                        not null,
				salt       bytea                                               not null,
				avatar     varchar(255),
				created_at timestamp with time zone default now()              not null,
				updated_at timestamp with time zone default now()              not null
			);
			
			create index if not exists idx_users_created_at
				on users (created_at);
			
			create index if not exists idx_users_first_name
				on users (first_name);
			
			create index if not exists idx_users_last_name
				on users (last_name);
			
			create unique index if not exists uix_users_email
				on users (email);
			
			create table roles
			(
				id   uuid default uuid_generate_v4() not null
					constraint roles_pk
						primary key,
				name varchar(32)                     not null
			);
			
			create unique index roles_name_uindex
				on roles (name);
			
			create table user_roles
			(
				id      uuid default uuid_generate_v4() not null
					constraint user_roles_pk
						primary key,
				user_id uuid                            not null
					constraint user_roles_users_id_fk
						references users
						on update cascade on delete cascade
						deferrable,
				role_id uuid                            not null
					constraint user_roles_roles_id_fk
						references roles
						on update restrict on delete restrict
						deferrable
			);
			
			create unique index user_roles_role_id_user_id_uindex
				on user_roles (role_id, user_id);
			
			INSERT INTO roles (name)
			VALUES ('account');
		`,
	},
}
