import {parseRatio, Ratio, ratioToFloat} from "./Ratio";
import {Flash, parseFlash} from "./Flash";
import {ExposureMode, parseExposureMode} from "./ExposureMode";
import {ExposureProgram, parseExposureProgram} from "./ExposureProgram";
import {LightSource, parseLightSource} from "./LightSource";
import {MeteringMode, parseMeteringMode} from "./MeteringMode";
import {parseWhiteBalance, WhiteBalance} from "./WhiteBalance";
import {parseSceneMode, SceneMode} from "./SceneMode";
import {ContrastProcessing, parseContrastProcessing} from "./ContrastProcessing";
import {parseSharpnessProcessing, SharpnessProcessing} from "./SharpnessProcessing";
import {parseSubjectDistanceRange, SubjectDistanceRange} from "./SubjectDistanceRange";

export interface ImageMetadata {
    aspectRatio?: Ratio,
    make?: string,
    model?: string,
    software?: string,
    copyright?: string,
    dateTimeCreated?: Date,
    dateTimeDigitized?: Date,
    dateTimeOriginal?: Date,
    digitalZoomRatio?: Ratio,
    exposure?: Ratio,
    exposureMode?: ExposureMode,
    exposureProgram?: ExposureProgram,
    shutterSpeed?: Ratio,
    aperture?: Ratio,
    flash?: Flash,
    focalLength?: Ratio,
    focalLength35mm?: Ratio,
    isoSpeedRating?: number,
    lightSource?: LightSource,
    meteringMode?: MeteringMode,
    whiteBalance?: WhiteBalance,
    sceneMode?: SceneMode,
    contrast?: ContrastProcessing,
    sharpness?: SharpnessProcessing,
    subjectDistance?: Ratio,
    subjectDistanceRange?: SubjectDistanceRange,
}

export function parseMetadata(metadata?: { [key: string]: string | undefined }): ImageMetadata | undefined {
    if (metadata === undefined) {
        return undefined;
    }
    return {
        aspectRatio: parseRatio(metadata["AspectRatio"]),
        make: metadata["Make"],
        model: metadata["Model"],
        software: metadata["Software"],
        copyright: metadata["Copyright"],
        dateTimeCreated: parseDate(metadata["DateTime"]),
        dateTimeDigitized: parseDate(metadata["DateTimeDigitized"]),
        dateTimeOriginal: parseDate(metadata["DateTimeOriginal"]),
        digitalZoomRatio: parseRatio(metadata["DigitalZoomRatio"]),
        exposure: parseRatio(metadata["ExposureBiasValue"]),
        exposureMode: parseExposureMode(metadata["ExposureMode"]),
        exposureProgram: parseExposureProgram(metadata["ExposureProgram"]),
        shutterSpeed: parseRatio(metadata["ExposureTime"]),
        aperture: parseRatio(metadata["FNumber"]),
        flash: parseFlash(metadata["Flash"], metadata["FlashEnergy"]),
        focalLength: parseRatio(metadata["FocalLength"]),
        focalLength35mm: parseRatio(metadata["FocalLengthIn35mmFilm"]),
        isoSpeedRating: parseNumber(metadata["ISOSpeedRatings"]),
        lightSource: parseLightSource(metadata["LightSource"]),
        meteringMode: parseMeteringMode(metadata["MeteringMode"]),
        whiteBalance: parseWhiteBalance(metadata["WhiteBalance"]),
        sceneMode: parseSceneMode(metadata["SceneMode"]),
        contrast: parseContrastProcessing(metadata["Contrast"]),
        sharpness: parseSharpnessProcessing(metadata["Sharpness"]),
        subjectDistance: parseSubjectDistance(metadata["SubjectDistance"]),
        subjectDistanceRange: parseSubjectDistanceRange(metadata["SubjectDistanceRange"]),
    }
}

export function parseDate(value: string | undefined): Date | undefined {
    if (value === undefined) {
        return undefined;
    }
    const split = value.split(" ");
    if (split.length !== 2) {
        return undefined;
    }
    const [date, time] = split;
    try {
        const parsed = new Date(
            date.replaceAll(":", "-") + " " + time
        );
        parsed.toISOString();
        return parsed;
    } catch (e) {
        return undefined;
    }
}

export function parseNumber(value: string | undefined): number | undefined {
    if (value === undefined) {
        return undefined;
    }
    const number = parseInt(value);
    if (isNaN(number) || number === Infinity || number === -Infinity) {
        return undefined;
    }
    return number;
}

export function ratioToTime(value: Ratio | undefined): string | undefined {
    if (value === undefined) {
        return undefined;
    }
    if (value.numerator > value.denominator) {
        return (value.numerator / value.denominator).toFixed(0) + "s";
    } else {
        return "1/" + (value.denominator / value.numerator).toFixed(0) + "s";
    }
}

export function ratioToDistance(value: Ratio | undefined): string | undefined {
    if (value === undefined) {
        return undefined;
    }
    if (value.numerator === Infinity) {
        return "∞";
    }
    const numeric = value.numerator / value.denominator;
    const magnitude = Math.floor(Math.log10(numeric));
    switch (magnitude) {
        case 0:
        case -1:
        case -2:
            return `${(numeric * 100).toFixed(0)} cm`;
        case -3:
            return `${(numeric * 1000).toFixed(0)} mm`;
        default:
            if (magnitude >= 0) return `${numeric.toFixed(0)} m`;
            else return `${numeric} m`;
    }
}

export function parseSubjectDistance(value: string | undefined): Ratio | undefined {
    const numericValue = parseRatio(value);
    if (numericValue === undefined) {
        return undefined;
    }
    if (ratioToFloat(numericValue) === 4294967295) {
        return {
            numerator: Infinity,
            denominator: 1,
        };
    }
    return numericValue;
}
