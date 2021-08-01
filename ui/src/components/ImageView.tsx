import {Image} from "../api/model/Image";
import React, {useMemo, useState} from "react";
import {useUpdateImage} from "../api/useUpdateImage";
import {useDeleteImage} from "../api/useDeleteImage";
import {parseMetadata} from "../metadata/ImageMetadata";
import ImageMetadataView from "./ImageMetadataView";
import {Button, List, ListItem, ListItemIcon, ListItemText, TextField} from "@material-ui/core";
import {Event, Info} from "@material-ui/icons";
import {File, Tag} from "mdi-material-ui";

export interface ImageProps {
    image: Image
}

export default function ImageView({image}: ImageProps) {
    const {mutate: update, error: updateError, isLoading: updateLoading} = useUpdateImage();
    const {mutate: remove, error: removeError, isLoading: removeLoading} = useDeleteImage();
    const [title, setTitle] = useState<string>(image.title);
    const [description, setDescription] = useState<string>(image.description);

    const metadata = useMemo(() =>
            parseMetadata(image.metadata),
        [image]);

    return (
        <div>
            <p>UpdateError: {JSON.stringify(updateError, null, 2)}</p>
            <p>RemoveError: {JSON.stringify(removeError, null, 2)}</p>
            <p>UpdateLoading: {JSON.stringify(updateLoading, null, 2)}</p>
            <p>RemoveLoading: {JSON.stringify(removeLoading, null, 2)}</p>
            <img src={image.url + "l"} alt=""/>
            <List dense>
                <ListItem dense>
                    <ListItemIcon><Info/></ListItemIcon>
                    <ListItemText primary="Id" secondary={image.id}/>
                </ListItem>
                <ListItem dense>
                    <ListItemText inset primary="Owner" secondary={image.owner}/>
                </ListItem>
                <ListItem dense>
                    <ListItemText inset primary="State" secondary={image.state}/>
                </ListItem>
                <ListItem dense>
                    <ListItemIcon><Tag/></ListItemIcon>
                    {/* TODO: Fix this ugly nesting */}
                    <ListItemText
                        primary="Title"
                        secondary={<TextField
                            fullWidth
                            value={title}
                            onChange={({target: {value}}) =>
                                setTitle(value)}
                        />}
                    />
                </ListItem>
                <ListItem dense>
                    {/* TODO: Fix this ugly nesting */}
                    <ListItemText
                        inset
                        primary="Description"
                        secondary={<TextField
                            fullWidth
                            value={description}
                            onChange={({target: {value}}) =>
                                setDescription(value)}
                        />}
                    />
                </ListItem>
                <ListItem dense>
                    <ListItemIcon><File/></ListItemIcon>
                    <ListItemText primary="Filename" secondary={image.original_name}/>
                </ListItem>
                <ListItem dense>
                    <ListItemText inset primary="MIME Type" secondary={image.mime_type}/>
                </ListItem>
                <ListItem dense>
                    <ListItemIcon><Event/></ListItemIcon>
                    <ListItemText primary="Uploaded At" secondary={image.created_at}/>
                </ListItem>
                <ListItem dense>
                    <ListItemText inset primary="Modified At" secondary={image.updated_at}/>
                </ListItem>
                <ImageMetadataView metadata={metadata}/>
            </List>
            <Button
                onClick={() => update({
                    ...image,
                    title,
                    description,
                })}
            >Save</Button>
            <Button
                onClick={() => remove(image)}
            >Delete</Button>
        </div>
    )
}
