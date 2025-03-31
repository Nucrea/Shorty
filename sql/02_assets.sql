create type assets_status as enum ('pending', 'created', 'deleted');

create table if not exists assets (
    id char(32) primary key,
    resource_id char(32) unique not null,
    size integer not null,
    hash char(128) not null,
    bucket varchar(256) not null,
    status assets_status not null default 'pending',
    -- report_count integer not null default 0,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now()
);

-- create index idx_assets_hash on assets using hash(hash);

