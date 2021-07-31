import {parseNumber} from "./ImageMetadata";

export enum SharpnessProcessing {
    NORMAL = 0,
    SOFT = 1,
    HARD = 2,
}

export function parseSharpnessProcessing(value: string | undefined): SharpnessProcessing | undefined {
    const numericValue = parseNumber(value);
    if (numericValue === undefined) {
        return undefined;
    }
    if (Object.values(SharpnessProcessing).includes(numericValue)) {
        return numericValue as SharpnessProcessing;
    }
    return undefined;
}
