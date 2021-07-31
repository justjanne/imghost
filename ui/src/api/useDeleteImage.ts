import {Image} from "./model/Image";
import {useBaseUrl} from "./baseUrlContext";
import {useMutation, useQueryClient} from "react-query";
import axios from "axios";

export const useDeleteImage = () => {
    const baseUrl = useBaseUrl();
    const queryClient = useQueryClient();
    return useMutation<unknown, unknown, Image>((image: Image) => axios.delete(
        `api/v1/images/${image.id}`,
        {
            baseURL: baseUrl
        }
    ), {
        onSuccess: () => {
            queryClient.invalidateQueries('image')
        },
    })
}
