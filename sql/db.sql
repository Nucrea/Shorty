create table if not exists shortlinks (
    id int primary key generated always as identity,
    short_id char(10) not null unique,
    url text not null,
    read_count int default 0,
    created_at timestamp not null default now()
);

create index if not exists idx_shortlinks_shortid on shortlinks(short_id);

create table if not exists images (
    id int primary key generated always as identity,
    short_id char(32) unique,
    size int not null,
    name varchar(256) not null,
    hash char(128) not null,
    image_id char(32) not null,
    thumbnail_id char(32) not null,
    created_at timestamp not null default now()
);

create index if not exists idx_image_short_id on images(short_id);

create index if not exists idx_image_hash on images(hash);