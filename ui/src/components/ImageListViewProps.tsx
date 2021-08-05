import {Fragment} from "react";
import {Link} from "react-router-dom";
import {Image} from "../api/model/Image";
import {parseRatio} from "../metadata/Ratio";

export interface ImageListViewProps {
    images: Image[],
}

export function ImageListView({images}: ImageListViewProps) {
    const rows: (Image | null)[][] = [];
    if (images) {
        let row = [];
        for (let i = 0; i < images.length; i++) {
            if (i % 4 === 0 && row.length !== 0) {
                rows.push(row);
                row = [];
            }
            row.push(images[i]);
        }
        rows.push([...row, ...([null, null, null].slice(0, 4 - row.length))]);
    }

    return (
        <Fragment>
            <style>
                {`img {
                    height: 100%;
                    max-width: 100%;
                    margin: 2px;
                }
                
                a:first-child img {
                margin-left: -8px;
                }
                
                a:last-child img {
                margin-right: -8px;
                }`}
            </style>
            <div style={{padding: "0 8px", fontSize: 0}}>
                {rows.map(row => {
                    let ratios = 0;
                    for (const image of row) {
                        const [w, h] = image ? determineAspectRatio(image) : [1, 1];
                        ratios += w / h;
                    }

                    return (
                        <div style={{aspectRatio: "" + ratios + "/" + 1}}>
                            {row.map(image => image && (
                                <Link to={`/i/${image.id}`}>
                                    <img src={image.url} alt=""/>
                                </Link>
                            ))}
                        </div>
                    )
                })}
            </div>
        </Fragment>
    )
}

function determineAspectRatio(image: Image): [number, number] {
    const aspectRatio = parseRatio(image.metadata["AspectRatio"]);
    if (aspectRatio === undefined) {
        console.log("Did not find aspect ratio: ", image.id);
        return [1, 1];
    } else {
        console.log("Found aspect ratio: ", image.id, aspectRatio.numerator, aspectRatio.denominator);
        return [aspectRatio.numerator, aspectRatio.denominator];
    }
}

