create table if not exists albums
(
    id          text            not null primary key,
    owner       text            not null,
    title       text default '' not null,
    description text default '' not null,
    created_at  timestamp       not null,
    updated_at  timestamp       not null
);

create index if not exists albums_owner_index
    on albums using hash (owner);

create index if not exists albums_created_at_index
    on albums (created_at);

create index if not exists albums_updated_at_index
    on albums (updated_at);
