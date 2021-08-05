import {useUpdateImage} from "../api/useUpdateImage";
import {useDeleteImage} from "../api/useDeleteImage";
import {useEffect, useMemo, useState} from "react";
import {parseMetadata} from "../metadata/ImageMetadata";
import {
    Button,
    CircularProgress,
    Grid,
    LinearProgress,
    List,
    ListItem,
    ListItemIcon,
    ListItemText,
    TextField
} from "@material-ui/core";
import {Delete, Event, Info, Save} from "@material-ui/icons";
import {File, Tag} from "mdi-material-ui";
import ImageMetadataView from "../components/ImageMetadataView";
import {useGetImage} from "../api/useGetImage";
import {useParams} from "react-router";
import {ErrorPortal} from "../components/ErrorContext";
import {ErrorAlert} from "../components/ErrorAlert";

export interface ImageDetailPageParams {
    imageId: string
}

export default function ImageDetailPage() {
    const {imageId} = useParams<ImageDetailPageParams>();
    const {data: image, error: imageError, isLoading: imageLoading} = useGetImage(imageId);
    const {mutate: update, error: updateError, isLoading: updateLoading} = useUpdateImage();
    const {mutate: remove, error: removeError, isLoading: removeLoading} = useDeleteImage();
    const [title, setTitle] = useState<string>(image?.title || "");
    const [description, setDescription] = useState<string>(image?.description || "");
    useEffect(() => setTitle(image?.title || ""), [image?.title]);
    useEffect(() => setDescription(image?.description || ""), [image?.description]);

    const metadata = useMemo(() => parseMetadata(image?.metadata), [image]);

    if (image === undefined || metadata === undefined) {
        return (
            <div>Error: 404</div>
        );
    }

    return (
        <div>
            {imageLoading && (
                <LinearProgress/>
            )}
            <ErrorPortal>
                <ErrorAlert severity="error" error={imageError}/>
                <ErrorAlert severity="error" error={updateError}/>
                <ErrorAlert severity="error" error={removeError}/>
            </ErrorPortal>
            <Grid container>
                <Grid item xs={12} md={6}>
                    <img src={image.url + "l"} alt={image.title} style={{width: "100%"}}/>
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
                    </List>
                    <Button
                        variant="contained"
                        color="primary"
                        disabled={updateLoading}
                        startIcon={updateLoading ? <CircularProgress style={{color: "#fff"}} size="1em"/> : <Save/>}
                        onClick={() => update({
                            ...image,
                            title,
                            description,
                        })}
                    >Save</Button>
                    <Button
                        variant="contained"
                        color="secondary"
                        disabled={removeLoading}
                        startIcon={removeLoading ? <CircularProgress style={{color: "#fff"}} size="1em"/> : <Delete/>}
                        onClick={() => remove(image)}
                    >Delete</Button>
                </Grid>
                <Grid item xs={12} md={6}>
                    <List dense>
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
                </Grid>
            </Grid>
        </div>
    )
}
