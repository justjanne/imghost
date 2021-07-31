import {parseNumber} from "./ImageMetadata";

export enum ContrastProcessing {
    NORMAL = 0,
    SOFT = 1,
    HARD = 2,
}

export function parseContrastProcessing(value: string | undefined): ContrastProcessing | undefined {
    const numericValue = parseNumber(value);
    if (numericValue === undefined) {
        return undefined;
    }
    if (Object.values(ContrastProcessing).includes(numericValue)) {
        return numericValue as ContrastProcessing;
    }
    return undefined;
}
