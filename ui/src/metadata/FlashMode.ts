export enum FlashMode {
    ALWAYS_ON = 1,
    ALWAYS_OFF = 2,
    AUTO = 3
}

export function parseFlashMode(value: number): FlashMode | undefined {
    if (Object.values(FlashMode).includes(value)) {
        return value as FlashMode;
    } else {
        return undefined;
    }
}
