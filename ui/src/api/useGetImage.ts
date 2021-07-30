import {useQuery} from "react-query";
import axios from "axios";
import {useBaseUrl} from "./baseUrlContext";
import {Image} from "./model/Image";

export const useGetImage = (imageId: string) => {
    const baseUrl = useBaseUrl();
    return useQuery(
        ["image", imageId],
        () => axios.get<Image>(
            `api/v1/images/${imageId}`,
            {
                baseURL: baseUrl
            }
        ).then(it => it.data),
        {
            keepPreviousData: true
        }
    );
}
