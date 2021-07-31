import {parseNumber} from "./ImageMetadata";

export enum SceneMode {
    STANDARD = 0,
    LANDSCAPE = 1,
    PORTRAIT = 2,
    NIGHT_SCENE = 3,
}

export function parseSceneMode(value: string | undefined): SceneMode | undefined {
    const numericValue = parseNumber(value);
    if (numericValue === undefined) {
        return undefined;
    }
    if (Object.values(SceneMode).includes(numericValue)) {
        return numericValue as SceneMode;
    }
    return undefined;
}
