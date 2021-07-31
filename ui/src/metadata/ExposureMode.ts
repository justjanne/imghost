import {parseNumber} from "./ImageMetadata";

export enum ExposureMode {
    AUTO = 0,
    MANUAL = 1,
    BRACKET = 2,
}

export function parseExposureMode(value: string | undefined): ExposureMode | undefined {
    const numericValue = parseNumber(value)
    if (numericValue === undefined) {
        return undefined;
    }
    if (Object.values(ExposureMode).includes(numericValue)) {
        return numericValue as ExposureMode;
    }
    return undefined;
}
