import {parseRatio, Ratio} from "./Ratio";
import {Flash} from "./Flash";
import {ExposureMode} from "./ExposureMode";
import {ExposureProgram} from "./ExposureProgram";
import {LightSource} from "./LightSource";
import {MeteringMode} from "./MeteringMode";
import {WhiteBalance} from "./WhiteBalance";
import {SceneMode} from "./SceneMode";
import {ContrastProcessing} from "./ContrastProcessing";
import {SharpnessProcessing} from "./SharpnessProcessing";
import {SubjectDistanceRange} from "./SubjectDistanceRange";

export interface ImageMetadata {
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
    exposureTime?: Ratio,
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
    subjectDistance?: number,
    subjectDistanceRange?: SubjectDistanceRange,
}

export function parseMetadata(metadata: { [key: string]: string }): ImageMetadata {
    return {
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
        exposureTime: parseRatio(metadata["ExposureTime"]),
        aperture: parseRatio(metadata["FNumber"]),
        flash: undefined,
        focalLength: parseRatio(metadata["FocalLength"]),
        focalLength35mm: parseRatio(metadata["FocalLengthIn35mmFilm"]),
        isoSpeedRating: parseNumber(metadata["ISOSpeedRatings"]),
        lightSource: parseLightSource(metadata["LightSource"]),
        meteringMode: parseMeteringMode(metadata["MeteringMode"]),
        whiteBalance: parseWhiteBalance(metadata["WhiteBalance"]),
        sceneMode: parseSceneMode(metadata["SceneMode"]),
        contrast: parseContrastProcessing(metadata["Contrast"]),
        sharpness: parseSharpnessProcessing(metadata["Sharpness"]),
        subjectDistance: parseNumber(metadata["SubjectDistance"]),
        subjectDistanceRange: parseSubjectDistanceRange(metadata["SubjectDistanceRange"]),
    }
}

export function parseDate(value: string): Date | undefined {
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

export function parseNumber(value: string): number | undefined {
    const number = parseInt(value);
    if (isNaN(number) || number === Infinity || number === -Infinity) {
        return undefined;
    }
    return number;
}

export function parseExposureMode(value: string): ExposureMode | undefined {
    const numericValue = parseNumber(value)
    if (numericValue === undefined) {
        return undefined;
    }
    if (numericValue in Object.values(ExposureMode)) {
        return numericValue as ExposureMode;
    }
    return undefined;
}

export function parseExposureProgram(value: string): ExposureProgram | undefined {
    const numericValue = parseNumber(value);
    if (numericValue === undefined) {
        return undefined;
    }
    if (numericValue in Object.values(ExposureProgram)) {
        return numericValue as ExposureProgram;
    }
    return undefined;
}

export function parseLightSource(value: string): LightSource | undefined {
    const numericValue = parseNumber(value);
    if (numericValue === undefined) {
        return undefined;
    }
    if (numericValue in Object.values(LightSource)) {
        return numericValue as LightSource;
    }
    return undefined;
}

export function parseMeteringMode(value: string): MeteringMode | undefined {
    const numericValue = parseNumber(value);
    if (numericValue === undefined) {
        return undefined;
    }
    if (numericValue in Object.values(MeteringMode)) {
        return numericValue as MeteringMode;
    }
    return undefined;
}

export function parseWhiteBalance(value: string): WhiteBalance | undefined {
    const numericValue = parseNumber(value);
    if (numericValue === undefined) {
        return undefined;
    }
    if (numericValue in Object.values(WhiteBalance)) {
        return numericValue as WhiteBalance;
    }
    return undefined;
}

export function parseSceneMode(value: string): SceneMode | undefined {
    const numericValue = parseNumber(value);
    if (numericValue === undefined) {
        return undefined;
    }
    if (numericValue in Object.values(SceneMode)) {
        return numericValue as SceneMode;
    }
    return undefined;
}

export function parseContrastProcessing(value: string): ContrastProcessing | undefined {
    const numericValue = parseNumber(value);
    if (numericValue === undefined) {
        return undefined;
    }
    if (numericValue in Object.values(ContrastProcessing)) {
        return numericValue as ContrastProcessing;
    }
    return undefined;
}

export function parseSharpnessProcessing(value: string): SharpnessProcessing | undefined {
    const numericValue = parseNumber(value);
    if (numericValue === undefined) {
        return undefined;
    }
    if (numericValue in Object.values(SharpnessProcessing)) {
        return numericValue as SharpnessProcessing;
    }
    return undefined;
}

export function parseSubjectDistanceRange(value: string): SubjectDistanceRange | undefined {
    const numericValue = parseNumber(value);
    if (numericValue === undefined) {
        return undefined;
    }
    if (numericValue in Object.values(SubjectDistanceRange)) {
        return numericValue as SubjectDistanceRange;
    }
    return undefined;
}

export function ratioToTime(value: Ratio | undefined): string | undefined {
    if (value === undefined) {
        return undefined;
    }
    if (value.numerator > value.denominator) {
        return (value.numerator / value.denominator).toFixed(0) + "s";
    } else {
        return "1/" + (value.denominator / value.numerator) + "s";
    }
}
