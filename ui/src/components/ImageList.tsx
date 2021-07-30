import React from 'react';
import '../App.css';
import {useListImages} from "../api/useListImages";
import ImageView from "./ImageView";

export default function ImageList() {
    const {status, data, error} = useListImages();
    return (
        <div>
            <p>{status}</p>
            <p>{error as string}</p>
            <ul>
                {data?.map(image => (
                    <ImageView
                        key={image.id}
                        image={image}
                    />
                ))}
            </ul>
        </div>
    );
}
