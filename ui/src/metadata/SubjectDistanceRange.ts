import {parseNumber} from "./ImageMetadata";

export enum SubjectDistanceRange {
    MACRO = 1,
    CLOSE = 2,
    DISTANT = 3,
}

export function parseSubjectDistanceRange(value: string | undefined): SubjectDistanceRange | undefined {
    const numericValue = parseNumber(value);
    if (numericValue === undefined) {
        return undefined;
    }
    if (Object.values(SubjectDistanceRange).includes(numericValue)) {
        return numericValue as SubjectDistanceRange;
    }
    return undefined;
}
