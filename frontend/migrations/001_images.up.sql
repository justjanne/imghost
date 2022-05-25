create table if not exists images
(
    id            text               not null primary key,
    owner         text               not null,
    title         text  default ''   not null,
    description   text  default ''   not null,
    original_name text  default ''   not null,
    type          text               not null,
    metadata      jsonb default '{}' not null,
    created_at    timestamp          not null,
    updated_at    timestamp          not null
);

create index if not exists images_owner_index
    on images using hash (owner);

create index if not exists images_created_at_index
    on images (created_at);

create index if not exists images_updated_at_index
    on images (updated_at);
