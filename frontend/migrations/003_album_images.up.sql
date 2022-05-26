create table if not exists album_images
(
    album       text            not null
        constraint album_images_albums_id_fk
            references albums
            on update cascade on delete cascade,
    image       text            not null
        constraint album_images_images_id_fk
            references images
            on update cascade on delete cascade,
    title       text default '' not null,
    description text default '' not null,
    position    integer         not null,
    constraint album_images_image_album_pk
        primary key (image, album)
);
