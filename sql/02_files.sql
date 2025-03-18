create table if not exists resources (
    id char(32) primary key,
    size integer not null,
    hash char(128) not null,
    active boolean not null default true,
    report_count integer not null default 0,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now()
);

create table if not exists images (
    id char(32) primary key,
    image_id references resources(id) not null,
    thumbnail_id references resources(id) not null,
    name varchar(256) not null,
    read_count integer not null default 0,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now()
);

create table if not exists files (
    id char(32) primary key,
    file_id references resources(id) not null,
    name varchar(256) not null,
    read_count integer not null default 0,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now()
);

