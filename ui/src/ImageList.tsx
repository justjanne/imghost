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
                    <li key={info.id}>
                        <p>{info.id}</p>
                        <p>{info.owner}</p>
                        <p>{info.title}</p>
                        <p>{info.description}</p>
                        <p>{info.original_name}</p>
                        <p>{info.mime_type}</p>
                        <p>{info.created_at}</p>
                        <p>{info.updated_at}</p>
                        <p>{info.state}</p>
                        <p>{info.url}</p>
                        <img src={info.url+"t"} alt=""/>
                    </li>
                ))}
            </ul>
        </div>
    );
}
