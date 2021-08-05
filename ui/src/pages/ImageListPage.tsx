import {Fragment} from "react";
import {useListImages} from "../api/useListImages";
import {ImageList, ImageListItem, ImageListItemBar, LinearProgress} from "@material-ui/core";
import {Link} from "react-router-dom";
import {ErrorPortal} from "../components/ErrorContext";
import {ErrorAlert} from "../components/ErrorAlert";

export default function ImageListPage() {
    const {data: images, error, isLoading} = useListImages();

    return (
        <Fragment>
            {isLoading && (
                <LinearProgress/>
            )}
            <ErrorPortal>
                <ErrorAlert severity="error" error={error}/>
            </ErrorPortal>
            <ImageList cols={5}>
                {images?.map(image => (
                    <ImageListItem component={Link} to={`/i/${image.id}`}>
                        <img src={image.url} alt={image.title}/>
                        <ImageListItemBar
                            title={image.title}
                            subtitle={image.original_name}
                        />
                    </ImageListItem>
                ))}
            </ImageList>
        </Fragment>
    );
}
