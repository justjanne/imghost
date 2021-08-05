import {Paper} from "@material-ui/core";

export interface ErrorAlertProps {
    severity: string,
    error: unknown,
}

export function ErrorAlert({severity, error}: ErrorAlertProps) {
    if (!error) {
        return null;
    }

    return (
        <Paper variant="outlined" color={severity}>
            <strong>Error</strong>: {"" + error}
        </Paper>
    )
}
