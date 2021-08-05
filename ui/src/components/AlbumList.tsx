import '../App.css';
import {useListAlbums} from "../api/useListAlbums";
import AlbumView from "./AlbumView";

export default function AlbumList() {
    const {status, data, error} = useListAlbums();
    return (
        <div>
            <p>{status}</p>
            <p>{error as string}</p>
            <ul>
                {data?.map(album => (
                    <AlbumView
                        key={album.id}
                        album={album}
                    />
                ))}
            </ul>
        </div>
    );
}
