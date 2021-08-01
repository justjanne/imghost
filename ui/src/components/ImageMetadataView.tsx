import React, {Fragment} from "react";
import {ImageMetadata, ratioToTime} from "../metadata/ImageMetadata";
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
import {ListItem, ListItemIcon, ListItemText} from "@material-ui/core";

export interface ImageMetadataViewProps {
    metadata: ImageMetadata
}

export default function ImageMetadataView({metadata}: ImageMetadataViewProps) {
    return (
        <Fragment>
            {metadata.make !== undefined && (
                <ListItem dense>
                    <ListItemIcon><PhotoCamera/></ListItemIcon>
                    <ListItemText primary="Make" secondary={metadata.make}/>
                </ListItem>
            )}
            {metadata.model !== undefined && (
                <ListItem dense>
                    <ListItemText inset primary="Model" secondary={metadata.model}/>
                </ListItem>
            )}
            {metadata.software !== undefined && (
                <ListItem dense>
                    <ListItemText inset primary="Software" secondary={metadata.software}/>
                </ListItem>
            )}
            {metadata.copyright !== undefined && (
                <ListItem dense>
                    <ListItemIcon><Copyright/></ListItemIcon>
                    <ListItemText primary="Copyright" secondary={metadata.copyright}/>
                </ListItem>
            )}
            {metadata.dateTimeCreated !== undefined && (
                <ListItem dense>
                    <ListItemIcon><Event/></ListItemIcon>
                    <ListItemText primary="Created At" secondary={metadata.dateTimeCreated.toISOString()}/>
                </ListItem>
            )}
            {metadata.dateTimeDigitized !== undefined && (
                <ListItem dense>
                    <ListItemText inset primary="Digitized At"
                                  secondary={metadata.dateTimeDigitized.toISOString()}/>
                </ListItem>
            )}
            {metadata.dateTimeOriginal !== undefined && (
                <ListItem dense>
                    <ListItemText inset primary="Shot At" secondary={metadata.dateTimeOriginal.toISOString()}/>
                </ListItem>
            )}
            {metadata.digitalZoomRatio !== undefined && (
                <ListItem dense>
                    <ListItemIcon><ZoomIn/></ListItemIcon>
                    <ListItemText primary="Zoom" secondary={`${ratioToFloat(metadata.digitalZoomRatio)}x`}/>
                </ListItem>
            )}
            {metadata.focalLength !== undefined && (
                <ListItem dense>
                    <ListItemIcon><AngleAcute/></ListItemIcon>
                    <ListItemText primary="Focal Length" secondary={`${ratioToFloat(metadata.focalLength)}mm`}/>
                </ListItem>
            )}
            {metadata.focalLength35mm !== undefined && (
                <ListItem dense>
                    <ListItemText inset primary="35mm equivalent"
                                  secondary={`${ratioToFloat(metadata.focalLength35mm)}mm`}/>
                </ListItem>
            )}
            {metadata.shutterSpeed !== undefined && (
                <ListItem dense>
                    <ListItemIcon><CameraTimer/></ListItemIcon>
                    <ListItemText primary="Shutter Speed" secondary={ratioToTime(metadata.shutterSpeed)}/>
                </ListItem>
            )}
            {metadata.aperture !== undefined && (
                <ListItem dense>
                    <ListItemIcon><Camera/></ListItemIcon>
                    <ListItemText primary="Aperture" secondary={ratioToFloat(metadata.aperture)}/>
                </ListItem>
            )}
            {metadata.isoSpeedRating !== undefined && (
                <ListItem dense>
                    <ListItemIcon>ISO</ListItemIcon>
                    <ListItemText primary="ISO" secondary={metadata.isoSpeedRating}/>
                </ListItem>
            )}
            {metadata.exposure !== undefined && (
                <ListItem dense>
                    <ListItemIcon><Exposure/></ListItemIcon>
                    <ListItemText primary="Exposure" secondary={ratioToFloat(metadata.exposure)}/>
                </ListItem>
            )}
            {metadata.exposureMode !== undefined && (
                <ListItem dense>
                    <ListItemText inset primary="Mode" secondary={ExposureMode[metadata.exposureMode]}/>
                </ListItem>
            )}
            {metadata.exposureProgram !== undefined && (
                <ListItem dense>
                    <ListItemText inset primary="Program" secondary={ExposureProgram[metadata.exposureProgram]}/>
                </ListItem>
            )}
            {metadata.meteringMode !== undefined && (
                <ListItem dense>
                    <ListItemText inset primary="Metering mode" secondary={MeteringMode[metadata.meteringMode]}/>
                </ListItem>
            )}
            {metadata.flash !== undefined && (
                <Fragment>
                    <ListItem dense>
                        <ListItemIcon><Flash/></ListItemIcon>
                        <ListItemText primary="Flash"
                                      secondary={metadata.flash.available ? "Available" : "Unavailable"}/>
                    </ListItem>
                    <ListItem dense>
                        <ListItemText inset primary="Fired" secondary={metadata.flash.fired ? "Yes" : "No"}/>
                    </ListItem>
                    <ListItem dense>
                        <ListItemText inset primary="Red Eye Reduction"
                                      secondary={metadata.flash.redEyeReduction ? "Yes" : "No"}/>
                    </ListItem>
                    {metadata.flash.mode !== undefined && (
                        <ListItem dense>
                            <ListItemText inset primary="Mode" secondary={FlashMode[metadata.flash.mode]}/>
                        </ListItem>
                    )}
                    {metadata.flash.strength !== undefined && (
                        <ListItem dense>
                            <ListItemText inset primary="Strength" secondary={`${metadata.flash.strength} BCPS`}/>
                        </ListItem>
                    )}
                    <ListItem dense>
                        <ListItemText inset primary="Strobe Detection"
                                      secondary={!metadata.flash.strobeDetection.available ? "Unvailable" :
                                          metadata.flash.strobeDetection.detected ? "Strobe detected" :
                                              "No strobe detected"}
                        />
                    </ListItem>
                </Fragment>
            )}
            {metadata.lightSource !== undefined && (
                <ListItem dense>
                    <ListItemIcon><WhiteBalanceIncandescent/></ListItemIcon>
                    <ListItemText primary="Light Source" secondary={LightSource[metadata.lightSource]}/>
                </ListItem>
            )}
            {metadata.whiteBalance !== undefined && (
                <ListItem dense>
                    <ListItemIcon><WhiteBalanceIncandescent/></ListItemIcon>
                    <ListItemText primary="White balance" secondary={WhiteBalance[metadata.whiteBalance]}/>
                </ListItem>
            )}
            {metadata.subjectDistance !== undefined && (
                <ListItem dense>
                    <ListItemIcon><ArrowExpandHorizontal/></ListItemIcon>
                    <ListItemText primary="Subject Distance" secondary={metadata.subjectDistance}/>
                </ListItem>
            )}
            {metadata.subjectDistanceRange !== undefined && (
                <ListItem dense>
                    <ListItemIcon><ArrowExpandHorizontal/></ListItemIcon>
                    <ListItemText primary="Subject Distance Range"
                                  secondary={SubjectDistanceRange[metadata.subjectDistanceRange]}/>
                </ListItem>
            )}
            {metadata.sceneMode !== undefined && (
                <p><b>Scene Mode</b>: {SceneMode[metadata.sceneMode]}</p>
            )}
            {metadata.contrast !== undefined && (
                <ListItem dense>
                    <ListItemIcon><BrightnessMedium/></ListItemIcon>
                    <ListItemText primary="Contrast Processing" secondary={ContrastProcessing[metadata.contrast]}/>
                </ListItem>
            )}
            {metadata.sharpness !== undefined && (
                <ListItem dense>
                    <ListItemIcon><Blur/></ListItemIcon>
                    <ListItemText primary="Sharpness Processing" secondary={SharpnessProcessing[metadata.sharpness]}/>
                </ListItem>
            )}
        </Fragment>
    );
}
