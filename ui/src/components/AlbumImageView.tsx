import React from "react";
import {AlbumImage} from "../api/model/AlbumImage";

export interface AlbumImageProps {
    image: AlbumImage
}

export default function AlbumImageView({image}: AlbumImageProps) {
    return (
        <div>
            <p>{image.image}</p>
            <p>{image.title}</p>
            <p>{image.description}</p>
            <img src={image.url + "t"} alt=""/>
        </div>
    )
}
