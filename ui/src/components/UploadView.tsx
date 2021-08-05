import {useUploadImage} from "../api/useUploadImage";
import {ErrorPortal} from "./ErrorContext";
import {ErrorAlert} from "./ErrorAlert";
import {LinearProgress} from "@material-ui/core";

export default function UploadView() {
    const {mutate: upload, error, isLoading} = useUploadImage();

    return (
        <div>
            {isLoading && (
                <LinearProgress/>
            )}
            <ErrorPortal>
                <ErrorAlert severity="error" error={error}/>
            </ErrorPortal>
            <input
                type="file"
                disabled={isLoading}
                onChange={async ({target}) => {
                    if (target.files) {
                        await upload(target.files)
                        target.files = null;
                    }
                }}
            />
        </div>
    )
}
