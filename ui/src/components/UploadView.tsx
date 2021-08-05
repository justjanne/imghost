import {useUploadImage} from "../api/useUploadImage";

export default function UploadView() {
    const {mutate: upload, error: uploadError, isLoading: uploadLoading} = useUploadImage();

    return (
        <div>
            <pre>Error: {JSON.stringify(uploadError, null, 2)}</pre>
            <pre>Loading: {JSON.stringify(uploadLoading, null, 2)}</pre>
            <input
                type="file"
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
