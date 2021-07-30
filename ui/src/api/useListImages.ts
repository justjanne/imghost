import {useQuery} from "react-query";
import axios from "axios";
import {useBaseUrl} from "./baseUrlContext";
import {Image} from "./model/Image";

export const useListImages = () => {
    const baseUrl = useBaseUrl();
    return useQuery(
        "images",
        () => axios.get<Image[]>(
            "api/v1/images",
            {
                baseURL: baseUrl
            }
        ).then(it => it.data),
        {
            keepPreviousData: true
        }
    );
}
