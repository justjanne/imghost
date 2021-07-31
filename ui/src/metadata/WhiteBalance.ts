import {parseNumber} from "./ImageMetadata";

export enum WhiteBalance {
    AUTO = 0,
    MANUAL = 1,
}

export function parseWhiteBalance(value: string | undefined): WhiteBalance | undefined {
    const numericValue = parseNumber(value);
    if (numericValue === undefined) {
        return undefined;
    }
    if (Object.values(WhiteBalance).includes(numericValue)) {
        return numericValue as WhiteBalance;
    }
    return undefined;
}
