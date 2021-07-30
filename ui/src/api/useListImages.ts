import {useQuery} from "react-query";
import axios from "axios";
import {useBaseUrl} from "./baseUrlContext";
import {ImageInfo} from "./model/ImageInfo";

export const useListImages = () => {
    const baseUrl = useBaseUrl();
    return useQuery(
        "connector-deployments",
        () => axios.get<ImageInfo[]>(
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
