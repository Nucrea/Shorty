create table if not exists shortlinks (
    id char(10) primary key,
    user_id int references users(id) default null,
    url text not null,
    read_count int default 0,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now()
);