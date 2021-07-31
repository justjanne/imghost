import {parseFlash} from "./Flash";
import {FlashMode} from "./FlashMode";

describe("parseFlash", () => {
    test("Flash did not fire", () => {
        expect(parseFlash(0x00.toString(), undefined)).toStrictEqual({
            available: true,
            fired: false,
            strobeDetection: {
                available: false,
                detected: false,
            },
            mode: undefined,
            redEyeReduction: false,
            strength: undefined,
        })
    })
    test("Flash fired", () => {
        expect(parseFlash(0x01.toString(), undefined)).toStrictEqual({
            available: true,
            fired: true,
            strobeDetection: {
                available: false,
                detected: false,
            },
            mode: undefined,
            redEyeReduction: false,
            strength: undefined,
        })
    })
    test("Strobe return light not detected", () => {
        expect(parseFlash(0x05.toString(), undefined)).toStrictEqual({
            available: true,
            fired: true,
            strobeDetection: {
                available: true,
                detected: false,
            },
            mode: undefined,
            redEyeReduction: false,
            strength: undefined,
        })
    })
    test("Strobe return light detected", () => {
        expect(parseFlash(0x07.toString(), undefined)).toStrictEqual({
            available: true,
            fired: true,
            strobeDetection: {
                available: true,
                detected: true,
            },
            mode: undefined,
            redEyeReduction: false,
            strength: undefined,
        })
    })
    test("Flash fired, compulsory flash mode", () => {
        expect(parseFlash(0x09.toString(), undefined)).toStrictEqual({
            available: true,
            fired: true,
            strobeDetection: {
                available: false,
                detected: false,
            },
            mode: FlashMode.ALWAYS_ON,
            redEyeReduction: false,
            strength: undefined,
        })
    })
    test("Flash fired, compulsory flash mode, return light not detected", () => {
        expect(parseFlash(0x0D.toString(), undefined)).toStrictEqual({
            available: true,
            fired: true,
            strobeDetection: {
                available: true,
                detected: false,
            },
            mode: FlashMode.ALWAYS_ON,
            redEyeReduction: false,
            strength: undefined,
        })
    })
    test("Flash fired, compulsory flash mode, return light detected", () => {
        expect(parseFlash(0x0F.toString(), undefined)).toStrictEqual({
            available: true,
            fired: true,
            strobeDetection: {
                available: true,
                detected: true,
            },
            mode: FlashMode.ALWAYS_ON,
            redEyeReduction: false,
            strength: undefined,
        })
    })
    test("Flash did not fire, compulsory flash mode", () => {
        expect(parseFlash(0x10.toString(), undefined)).toStrictEqual({
            available: true,
            fired: false,
            strobeDetection: {
                available: false,
                detected: false,
            },
            mode: FlashMode.ALWAYS_OFF,
            redEyeReduction: false,
            strength: undefined,
        })
    })
    test("Flash did not fire, auto mode", () => {
        expect(parseFlash(0x18.toString(), undefined)).toStrictEqual({
            available: true,
            fired: false,
            strobeDetection: {
                available: false,
                detected: false,
            },
            mode: FlashMode.AUTO,
            redEyeReduction: false,
            strength: undefined,
        })
    })
    test("Flash fired, auto mode", () => {
        expect(parseFlash(0x19.toString(), undefined)).toStrictEqual({
            available: true,
            fired: true,
            strobeDetection: {
                available: false,
                detected: false,
            },
            mode: FlashMode.AUTO,
            redEyeReduction: false,
            strength: undefined,
        })
    })
    test("Flash fired, auto mode, return light not detected", () => {
        expect(parseFlash(0x1D.toString(), undefined)).toStrictEqual({
            available: true,
            fired: true,
            strobeDetection: {
                available: true,
                detected: false,
            },
            mode: FlashMode.AUTO,
            redEyeReduction: false,
            strength: undefined,
        })
    })
    test("Flash fired, auto mode, return light detected", () => {
        expect(parseFlash(0x1F.toString(), undefined)).toStrictEqual({
            available: true,
            fired: true,
            strobeDetection: {
                available: true,
                detected: true,
            },
            mode: FlashMode.AUTO,
            redEyeReduction: false,
            strength: undefined,
        })
    })
    test("No flash function", () => {
        expect(parseFlash(0x20.toString(), undefined)).toStrictEqual({
            available: false,
            fired: false,
            strobeDetection: {
                available: false,
                detected: false,
            },
            mode: undefined,
            redEyeReduction: false,
            strength: undefined,
        })
    })
    test("Flash fired, red-eye reduction mode", () => {
        expect(parseFlash(0x41.toString(), undefined)).toStrictEqual({
            available: true,
            fired: true,
            strobeDetection: {
                available: false,
                detected: false,
            },
            mode: undefined,
            redEyeReduction: true,
            strength: undefined,
        })
    })
    test("Flash fired, red-eye reduction mode, return light not detected", () => {
        expect(parseFlash(0x45.toString(), undefined)).toStrictEqual({
            available: true,
            fired: true,
            strobeDetection: {
                available: true,
                detected: false,
            },
            mode: undefined,
            redEyeReduction: true,
            strength: undefined,
        })
    })
    test("Flash fired, red-eye reduction mode, return light detected", () => {
        expect(parseFlash(0x47.toString(), undefined)).toStrictEqual({
            available: true,
            fired: true,
            strobeDetection: {
                available: true,
                detected: true,
            },
            mode: undefined,
            redEyeReduction: true,
            strength: undefined,
        })
    })
    test("Flash fired, compulsory flash mode, red-eye reduction mode", () => {
        expect(parseFlash(0x49.toString(), undefined)).toStrictEqual({
            available: true,
            fired: true,
            strobeDetection: {
                available: false,
                detected: false,
            },
            mode: FlashMode.ALWAYS_ON,
            redEyeReduction: true,
            strength: undefined,
        })
    })
    test("Flash fired, compulsory flash mode, red-eye reduction mode, return light not detected", () => {
        expect(parseFlash(0x4D.toString(), undefined)).toStrictEqual({
            available: true,
            fired: true,
            strobeDetection: {
                available: true,
                detected: false,
            },
            mode: FlashMode.ALWAYS_ON,
            redEyeReduction: true,
            strength: undefined,
        })
    })
    test("Flash fired, compulsory flash mode, red-eye reduction mode, return light detected", () => {
        expect(parseFlash(0x4F.toString(), undefined)).toStrictEqual({
            available: true,
            fired: true,
            strobeDetection: {
                available: true,
                detected: true,
            },
            mode: FlashMode.ALWAYS_ON,
            redEyeReduction: true,
            strength: undefined,
        })
    })
    test("Flash fired, auto mode, red-eye reduction mode", () => {
        expect(parseFlash(0x59.toString(), undefined)).toStrictEqual({
            available: true,
            fired: true,
            strobeDetection: {
                available: false,
                detected: false,
            },
            mode: FlashMode.AUTO,
            redEyeReduction: true,
            strength: undefined,
        })
    })
    test("Flash fired, auto mode, return light not detected, red-eye reduction mode", () => {
        expect(parseFlash(0x5D.toString(), undefined)).toStrictEqual({
            available: true,
            fired: true,
            strobeDetection: {
                available: true,
                detected: false,
            },
            mode: FlashMode.AUTO,
            redEyeReduction: true,
            strength: undefined,
        })
    })
    test("Flash fired, auto mode, return light detected, red-eye reduction mode", () => {
        expect(parseFlash(0x5F.toString(), undefined)).toStrictEqual({
            available: true,
            fired: true,
            strobeDetection: {
                available: true,
                detected: true,
            },
            mode: FlashMode.AUTO,
            redEyeReduction: true,
            strength: undefined,
        })
    })
})
