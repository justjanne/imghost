import AlbumImageView from "./AlbumImageView";
import React from "react";
import {Album} from "../api/model/Album";

export interface AlbumProps {
    album: Album
}

export default function AlbumView({album}: AlbumProps) {
    return (
        <div>
            <p>{album.id}</p>
            <p>{album.owner}</p>
            <p>{album.title}</p>
            <p>{album.description}</p>
            <p>{album.created_at}</p>
            <p>{album.updated_at}</p>
            <ul>
                {album.images.map(image => (
                    <AlbumImageView
                        key={image.image}
                        image={image}
                    />
                ))}
            </ul>
        </div>
    );
}
