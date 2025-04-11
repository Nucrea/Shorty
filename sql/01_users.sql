create table if not exists users (
    id int primary key generated always as identity,
    email varchar(256) unique not null,
    secret varchar(256) not null,
    verified boolean not null default true,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now()
);

create index idx_users_email on users(email);