import {FlashMode, parseFlashMode} from "./FlashMode";
import {parseRatio, Ratio} from "./Ratio";
import {parseNumber} from "./ImageMetadata";

export interface Flash {
    available: boolean,
    fired: boolean,
    strobeDetection: {
        available: boolean,
        detected: boolean,
    },
    mode: FlashMode | undefined,
    redEyeReduction: boolean,
    strength: Ratio | undefined,
}

const MASK_FIRED = 0x0001;
const MASK_STROBE_DETECTED = 0x0002;
const MASK_STROBE_DETECTION_AVAILABLE = 0x0004;
const MASK_MODE = 0x0003;
const MASK_UNAVAILABLE = 0x0020;
const MASK_REDEYE = 0x0040;

export function parseFlash(flash: string | undefined, strength: string | undefined): Flash | undefined {
    const value = parseNumber(flash);
    if (value === undefined) {
        return undefined;
    }

    return {
        available: (value & MASK_UNAVAILABLE) === 0,
        fired: (value & MASK_FIRED) !== 0,
        strobeDetection: {
            available: (value & MASK_STROBE_DETECTION_AVAILABLE) !== 0,
            detected: (value & MASK_STROBE_DETECTED) !== 0,
        },
        redEyeReduction: (value & MASK_REDEYE) !== 0,
        strength: parseRatio(strength),
        mode: parseFlashMode((value >> 3) & MASK_MODE),
    };
}
