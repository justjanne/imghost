import {FlashMode} from "./FlashMode";

export interface Flash {
    available: boolean,
    fired: boolean,
    strobeDetection: {
        available: boolean,
        detected: boolean,
    },
    mode: FlashMode | null,
    redEyeReduction: boolean,
}
