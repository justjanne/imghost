import {useQuery} from "react-query";
import axios from "axios";
import {useBaseUrl} from "./baseUrlContext";
import {Album} from "./model/Album";

export const useListAlbums = () => {
    const baseUrl = useBaseUrl();
    return useQuery(
        "albums",
        () => axios.get<Album[]>(
            "api/v1/albums",
            {
                baseURL: baseUrl
            }
        ).then(it => it.data),
        {
            keepPreviousData: true
        }
    );
}
