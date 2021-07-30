import React from 'react';
import './App.css';
import {useListImages} from "./api/useListImages";

export default function ImageList() {
    const {status, data, error} = useListImages();
    return (
        <div>
            <p>{status}</p>
            <p>{error as string}</p>
            <ul>
                {data?.map(info => (
                    <li>
                        <p>{info.image?.id}</p>
                        <p>{info.image?.title}</p>
                        <p>{info.image?.description}</p>
                        <p>{info.image?.original_name}</p>
                        <p>{info.image?.mime_type}</p>
                        <p>{info.image?.created_at}</p>
                        <p>{info.image?.updated_at}</p>
                        <p>{info.state}</p>
                        <img src={info.url} alt=""/>
                    </li>
                ))}
            </ul>
        </div>
    );
}
