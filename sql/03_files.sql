create table if not exists files (
    id int primary key generated always as identity,
    short_id char(32) unique,
    size int not null,
    name varchar(256) not null,
    resource_id char(32) unique,
    created_at timestamp not null default now()
);

create index if not exists idx_files_short_id on files(short_id);