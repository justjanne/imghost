create table if not exists albums
(
    id          text not null primary key,
    owner       text,
    title       text,
    description text,
    created_at  timestamp,
    updated_at  timestamp
);

create index if not exists albums_owner_index
    on albums using hash (owner);

create index if not exists albums_created_at_index
    on albums (created_at);

create index if not exists albums_updated_at_index
    on albums (updated_at);

