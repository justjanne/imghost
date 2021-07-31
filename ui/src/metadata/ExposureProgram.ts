import {parseNumber} from "./ImageMetadata";

export enum ExposureProgram {
    MANUAL = 1,
    NORMAL = 2,
    APERTURE_PRIORITY = 3,
    SHUTTER_PRIORITY = 4,
    CREATIVE = 5,
    ACTION = 6,
    PORTRAIT = 7,
    LANDSCAPE = 8,
}

export function parseExposureProgram(value: string | undefined): ExposureProgram | undefined {
    const numericValue = parseNumber(value);
    if (numericValue === undefined) {
        return undefined;
    }
    if (Object.values(ExposureProgram).includes(numericValue)) {
        return numericValue as ExposureProgram;
    }
    return undefined;
}
