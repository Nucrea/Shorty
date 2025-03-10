create table if not exists shortlinks (
    id int primary key generated always as identity,
    short_id varchar(10) not null unique,
    url text not null,
    read_count int default 0,
    created_at timestamp not null default now()
);

create index if not exists idx_shortlinks_shortid on shortlinks(short_id);

create table if not exists images (
    id int primary key generated always as identity,
    short_id varchar(32) unique,
    size int not null,
    name varchar(256) not null,
    image_id varchar(32) unique,
    thumbnail_id varchar(32) unique,
    created_at timestamp not null default now()
);

create index if not exists idx_image_short_id on images(short_id);