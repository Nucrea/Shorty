create table if not exists shortlinks (
    id int primary key generated always as identity,
    short_id char(10) not null unique,
    url text not null,
    read_count int default 0,
    created_at timestamp not null default now()
);

create index if not exists idx_shortlinks_shortid on shortlinks(short_id);