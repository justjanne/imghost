import {Image} from "./model/Image";
import {useBaseUrl} from "./baseUrlContext";
import {useMutation, useQueryClient} from "react-query";
import axios from "axios";

export const useUpdateImage = () => {
    const baseUrl = useBaseUrl();
    const queryClient = useQueryClient();
    return useMutation<void, unknown, Image>((image: Image) => axios.post(
        `api/v1/images/${image.id}`,
        image,
        {
            baseURL: baseUrl
        }
    ), {
        onSuccess: () => {
            queryClient.invalidateQueries('image')
        },
    })
}
