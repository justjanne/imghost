import {parseNumber} from "./ImageMetadata";

export enum LightSource {
    DAYLIGHT = 1,
    FLUORESCENT = 2,
    INCANDESCENT = 3,
    FLASH = 4,
    FINE_WEATHER = 9,
    CLOUDY_WEATHER = 10,
    SHADE = 11,
    FLUORESCENT_6400K = 12,
    FLUORESCENT_5000K = 13,
    FLUORESCENT_4200K = 14,
    FLUORESCENT_3450K = 15,
    STANDARD_LIGHT_A = 17,
    STANDARD_LIGHT_B = 18,
    STANDARD_LIGHT_C = 10,
    D55 = 20,
    D65 = 21,
    D75 = 22,
    D50 = 23,
    ISO_STUDIO_INCANDESCENT = 24,
}

export function parseLightSource(value: string | undefined): LightSource | undefined {
    const numericValue = parseNumber(value);
    if (numericValue === undefined) {
        return undefined;
    }
    if (Object.values(LightSource).includes(numericValue)) {
        return numericValue as LightSource;
    }
    return undefined;
}
