import {Image} from "../api/model/Image";
import React, {Fragment, useMemo, useState} from "react";
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
import {BrightnessMedium, Camera, Copyright, Event, Exposure, PhotoCamera, ZoomIn} from "@material-ui/icons";
import {AngleAcute, ArrowExpandHorizontal, Blur, CameraTimer, Flash, WhiteBalanceIncandescent} from "mdi-material-ui";
import {FlashMode} from "../metadata/FlashMode";

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
            {metadata.make !== undefined && (
                <p><b>Make</b>: {metadata.make}</p>
            )}
            {metadata.model !== undefined && (
                <p><PhotoCamera/><b>Model</b>: {metadata.model}</p>
            )}
            {metadata.software !== undefined && (
                <p><b>Software</b>: {metadata.software}</p>
            )}
            {metadata.copyright !== undefined && (
                <p><Copyright/><b>Copyright</b>: {metadata.copyright}</p>
            )}
            {metadata.dateTimeCreated !== undefined && (
                <p><Event/><b>DateTime Created</b>: {metadata.dateTimeCreated?.toISOString()}</p>
            )}
            {metadata.dateTimeDigitized !== undefined && (
                <p><Event/><b>DateTime Digitized</b>: {metadata.dateTimeDigitized?.toISOString()}</p>
            )}
            {metadata.dateTimeOriginal !== undefined && (
                <p><Event/><b>DateTime Original</b>: {metadata.dateTimeOriginal?.toISOString()}</p>
            )}
            {metadata.digitalZoomRatio !== undefined && (
                <p><ZoomIn/><b>Digital Zoom</b>: {ratioToFloat(metadata.digitalZoomRatio)}</p>
            )}
            {metadata.exposure !== undefined && (
                <p><Exposure/><b>Exposure</b>: {ratioToFloat(metadata.exposure)}</p>
            )}
            {metadata.exposureMode !== undefined && (
                <p><Exposure/><b>Exposure Mode</b>: {ExposureMode[metadata.exposureMode]}</p>
            )}
            {metadata.exposureProgram !== undefined && (
                <p><Exposure/><b>Exposure Program</b>: {ExposureProgram[metadata.exposureProgram]}</p>
            )}
            {metadata.shutterSpeed !== undefined && (
                <p><CameraTimer/><b>Shutter Speed</b>: {ratioToTime(metadata.shutterSpeed)}</p>
            )}
            {metadata.aperture !== undefined && (
                <p><Camera/><b>Aperture</b>: {ratioToFloat(metadata.aperture)}</p>
            )}
            {metadata.focalLength !== undefined && (
                <p><AngleAcute/><b>Focal Length</b>: {ratioToFloat(metadata.focalLength)}mm</p>
            )}
            {metadata.focalLength35mm !== undefined && (
                <p><AngleAcute/><b>Focal Length (35mm equivalent)</b>: {ratioToFloat(metadata.focalLength35mm)}mm</p>
            )}
            {metadata.isoSpeedRating !== undefined && (
                <p><b>ISO</b>: {metadata.isoSpeedRating}</p>
            )}
            {metadata.flash !== undefined && (
                <Fragment>
                    <p><Flash/><b>Flash</b></p>
                    <p><b>Available</b>: {metadata.flash.available ? "Yes" : "No"}</p>
                    <p><b>Fired</b>: {metadata.flash.fired ? "Yes" : "No"}</p>
                    <p><b>Red Eye Reduction</b>: {metadata.flash.redEyeReduction ? "Yes" : "No"}</p>
                    <p><b>Strobe Detection Available</b>: {metadata.flash.strobeDetection.available ? "Yes" : "No"}</p>
                    <p><b>Strobe Detection Used</b>: {metadata.flash.strobeDetection.detected ? "Yes" : "No"}</p>
                    {metadata.flash.mode !== undefined && (
                        <p><b>Flash Mode</b>: {FlashMode[metadata.flash.mode]}</p>
                    )}
                </Fragment>
            )}
            {metadata.lightSource !== undefined && (
                <p><WhiteBalanceIncandescent/><b>Light source</b>: {LightSource[metadata.lightSource]}</p>
            )}
            {metadata.meteringMode !== undefined && (
                <p><b>Metering mode</b>: {MeteringMode[metadata.meteringMode]}</p>
            )}
            {metadata.whiteBalance !== undefined && (
                <p><WhiteBalanceIncandescent/><b>White balance</b>: {WhiteBalance[metadata.whiteBalance]}</p>
            )}
            {metadata.sceneMode !== undefined && (
                <p><b>Scene Mode</b>: {SceneMode[metadata.sceneMode]}</p>
            )}
            {metadata.contrast !== undefined && (
                <p><BrightnessMedium/><b>Contrast Processing</b>: {ContrastProcessing[metadata.contrast]}</p>
            )}
            {metadata.sharpness !== undefined && (
                <p><Blur/><b>Sharpness Processing</b>: {SharpnessProcessing[metadata.sharpness]}</p>
            )}
            {metadata.subjectDistance !== undefined && (
                <p><ArrowExpandHorizontal/><b>Subject Distance</b>: {metadata.subjectDistance}</p>
            )}
            {metadata.subjectDistanceRange !== undefined && (
                <p><ArrowExpandHorizontal/><b>Subject Distance
                    Range</b>: {SubjectDistanceRange[metadata.subjectDistanceRange]}</p>
            )}
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
