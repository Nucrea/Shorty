create type assets_status as enum ('pending', 'created', 'deleted');

create table if not exists assets (
    id char(32) primary key,
    size integer not null,
    hash char(128) not null,
    bucket varchar(256) not null,
    status assets_status not null default 'pending',
    -- report_count integer not null default 0,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now()
);

create table if not exists images (
    id char(32) primary key,
    original_id char(32)references assets(id) not null,
    thumbnail_id char(32) references assets(id) not null,
    name varchar(256) not null,
    -- read_count integer not null default 0,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now()
);

create table if not exists files (
    id char(32) primary key,
    file_id char(32) references assets(id) not null,
    name varchar(256) not null,
    -- read_count integer not null default 0,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now()
);

