import {Image} from "../api/model/Image";
import React, {useMemo, useState} from "react";
import {useUpdateImage} from "../api/useUpdateImage";
import {useDeleteImage} from "../api/useDeleteImage";
import {parseMetadata, ratioToTime} from "../metadata/ImageMetadata";
import {ratioToFloat} from "../metadata/Ratio";
import {ExposureMode} from "../metadata/ExposureMode";
import {ExposureProgram} from "../metadata/ExposureProgram";
import {LightSource} from "../metadata/LightSource";
import {MeteringMode} from "../metadata/MeteringMode";
import {WhiteBalance} from "../metadata/WhiteBalance";
import {SceneMode} from "../metadata/SceneMode";
import {ContrastProcessing} from "../metadata/ContrastProcessing";
import {SharpnessProcessing} from "../metadata/SharpnessProcessing";
import {SubjectDistanceRange} from "../metadata/SubjectDistanceRange";

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
            <h3>Metadata</h3>
            <p><b>Make</b>: {metadata.make}</p>
            <p><b>Model</b>: {metadata.model}</p>
            <p><b>Software</b>: {metadata.software}</p>
            <p><b>Copyright</b>: {metadata.copyright}</p>
            <p><b>DateTime Created</b>: {metadata.dateTimeCreated?.toISOString()}</p>
            <p><b>DateTime Digitized</b>: {metadata.dateTimeDigitized?.toISOString()}</p>
            <p><b>DateTime Original</b>: {metadata.dateTimeOriginal?.toISOString()}</p>
            <p><b>Digital Zoom</b>: {ratioToFloat(metadata.digitalZoomRatio)}</p>
            <p><b>Exposure</b>: {ratioToFloat(metadata.exposure)}</p>
            <p><b>Exposure Mode</b>: {metadata.exposureMode !== undefined ?
                ExposureMode[metadata.exposureMode] : "null"}</p>
            <p><b>Exposure Program</b>: {metadata.exposureProgram !== undefined ?
                ExposureProgram[metadata.exposureProgram] : "null"}</p>
            <p><b>Exposure Time</b>: {ratioToTime(metadata.exposureTime)}</p>
            <p><b>Aperture</b>: {ratioToFloat(metadata.aperture)}</p>
            <p><b>Focal Length</b>: {ratioToFloat(metadata.focalLength)}</p>
            <p><b>Focal Length (35mm equivalent)</b>: {ratioToFloat(metadata.focalLength35mm)}</p>
            <p><b>ISO</b>: {metadata.isoSpeedRating}</p>
            <p><b>Light source</b>: {metadata.lightSource !== undefined ?
                LightSource[metadata.lightSource] : "null"}</p>
            <p><b>Metering mode</b>: {metadata.meteringMode !== undefined ?
                MeteringMode[metadata.meteringMode] : "null"}</p>
            <p><b>White balance</b>: {metadata.whiteBalance !== undefined ?
                WhiteBalance[metadata.whiteBalance] : "null"}</p>
            <p><b>Scene Mode</b>: {metadata.sceneMode !== undefined ?
                SceneMode[metadata.sceneMode] : "null"}</p>
            <p><b>Contrast Processing</b>: {metadata.contrast !== undefined ?
                ContrastProcessing[metadata.contrast] : "null"}</p>
            <p><b>Sharpness Processing</b>: {metadata.sharpness !== undefined ?
                SharpnessProcessing[metadata.sharpness] : "null"}</p>
            <p><b>Subject Distance</b>: {metadata.subjectDistance}</p>
            <p><b>Subject Distance Range</b>: {metadata.subjectDistanceRange !== undefined ?
                SubjectDistanceRange[metadata.subjectDistanceRange] : "null"}</p>
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
