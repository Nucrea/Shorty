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