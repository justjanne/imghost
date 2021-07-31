import {Image} from "../api/model/Image";
import React, {useState} from "react";
import {useUpdateImage} from "../api/useUpdateImage";
import {useDeleteImage} from "../api/useDeleteImage";

export interface ImageProps {
    image: Image
}

export default function ImageView({image}: ImageProps) {
    const {mutate: update, error: updateError, isLoading: updateLoading} = useUpdateImage();
    const {mutate: remove, error: removeError, isLoading: removeLoading} = useDeleteImage();
    const [title, setTitle] = useState<string>(image.title);
    const [description, setDescription] = useState<string>(image.description);

    return (
        <div>
            <p>UpdateError: {JSON.stringify(updateError, null, 2)}</p>
            <p>RemoveError: {JSON.stringify(removeError, null, 2)}</p>
            <p>UpdateLoading: {JSON.stringify(updateLoading, null, 2)}</p>
            <p>RemoveLoading: {JSON.stringify(removeLoading, null, 2)}</p>
            <p>{image.id}</p>
            <p>{image.owner}</p>
            <label>
                Title
                <input
                    type="text"
                    value={title}
                    onChange={({target: {value}}) =>
                        setTitle(value)}
                />
            </label>
            <br/>
            <label>
                Description
                <input
                    type="text"
                    value={description}
                    onChange={({target: {value}}) =>
                        setDescription(value)}
                />
            </label>
            <p>{image.original_name}</p>
            <p>{image.mime_type}</p>
            <p>{image.created_at}</p>
            <p>{image.updated_at}</p>
            <p>{image.state}</p>
            <img src={image.url + "t"} alt=""/>
            <br/>
            <input
                type="submit"
                value="Save"
                onClick={() => update({
                    ...image,
                    title,
                    description,
                })}
            />
            <input
                type="submit"
                value="Delete"
                onClick={() => remove(image)}
            />
        </div>
    )
}
