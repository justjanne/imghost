import {Image} from "../api/model/Image";
import React from "react";

export interface ImageProps {
    image: Image
}

export default function ImageView({image}: ImageProps) {
    return (
        <div>
            <p>{image.id}</p>
            <p>{image.owner}</p>
            <p>{image.title}</p>
            <p>{image.description}</p>
            <p>{image.original_name}</p>
            <p>{image.mime_type}</p>
            <p>{image.created_at}</p>
            <p>{image.updated_at}</p>
            <p>{image.state}</p>
            <img src={image.url + "t"} alt=""/>
        </div>
    )
}
