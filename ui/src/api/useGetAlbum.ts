import {useQuery} from "react-query";
import axios from "axios";
import {useBaseUrl} from "./baseUrlContext";
import {Album} from "./model/Album";

export const useGetAlbum = (albumId: string) => {
    const baseUrl = useBaseUrl();
    return useQuery(
        ["album", albumId],
        () => axios.get<Album>(
            `api/v1/albums/${albumId}`,
            {
                baseURL: baseUrl
            }
        ).then(it => it.data),
        {
            keepPreviousData: true
        }
    );
}
