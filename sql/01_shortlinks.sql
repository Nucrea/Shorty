create table if not exists shortlinks (
    id char(10) primary key,
    url text not null,
    read_count int default 0,
    created_at timestamp not null default now()
);