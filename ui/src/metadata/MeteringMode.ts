import {parseNumber} from "./ImageMetadata";

export enum MeteringMode {
    AVERAGE = 1,
    CENTER_WEIGHTED_AVERAGE = 2,
    SPOT = 3,
    MULTI_SPOT = 4,
    PATTERN = 5,
    PARTIAL = 6,
}

export function parseMeteringMode(value: string | undefined): MeteringMode | undefined {
    const numericValue = parseNumber(value);
    if (numericValue === undefined) {
        return undefined;
    }
    if (Object.values(MeteringMode).includes(numericValue)) {
        return numericValue as MeteringMode;
    }
    return undefined;
}
