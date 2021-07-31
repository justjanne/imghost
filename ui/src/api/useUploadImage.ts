import {Image} from "./model/Image";
import {useBaseUrl} from "./baseUrlContext";
import {useMutation, useQueryClient} from "react-query";
import axios from "axios";

export const useUploadImage = () => {
    const baseUrl = useBaseUrl();
    const queryClient = useQueryClient();
    return useMutation<Image, unknown, FileList>((files: FileList) => {
        const formData = new FormData();
        for (let i = 0; i < files.length; i++) {
            formData.append("images", files[i]);
        }
        return axios.post<Image>(
            `api/v1/images`,
            formData,
            {
                baseURL: baseUrl
            }
        ).then(it => it.data);
    }, {
        onSuccess: () => {
            queryClient.invalidateQueries('image')
        },
    })
}
