export interface Ratio {
    numerator: number,
    denominator: number
}

export function parseRatio(value: string | undefined): Ratio | undefined {
    if (value === undefined) {
        return undefined;
    }
    const splitValues = value.split("/");
    if (splitValues.length < 1) {
        return undefined;
    }
    const numerator = parseInt(splitValues[0]);
    if (isNaN(numerator)) {
        return undefined;
    }
    let denominator;
    if (splitValues.length === 1) {
        denominator = 1;
    } else {
        denominator = parseInt(splitValues[1]);
        if (isNaN(denominator) || denominator === 0) {
            return undefined;
        }
    }

    return {numerator, denominator};
}

export function ratioToFloat(ratio: Ratio | undefined): number | undefined {
    if (ratio === undefined) {
        return undefined;
    }
    return ratio.numerator / ratio.denominator;
}
