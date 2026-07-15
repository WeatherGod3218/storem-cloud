import {
  Card,
  CardContent,
  CardTitle,
} from "@/components/ui/card"

import { AspectRatio } from "@/components/ui/aspect-ratio"

type VideoData = {
    row_id: string,
    s3_id: string,

    custom_title: string | null,
    custom_description: string | null,

    username: string,
    filename: string,
    video_url: string,
}

export const VideoCard = (props: VideoData) => {
    return (
        <Card className="m-3">
            <CardContent>
                <AspectRatio ratio={16 / 9} className="w-3/4">
                    <video
                    className="w-full h-full rounded-md object-cover"
                    src={props.video_url}
                    controls
                    preload="metadata"
                    >  Your browser does not support the video tag.</video>
                </AspectRatio>
            </CardContent>
            <CardContent>
                <CardTitle>
                    {props.custom_title ? props.custom_title : props.filename}
                </CardTitle>
            </CardContent>
        </Card>
    )
}

export const SkeletonVideoCard = () => {
    return (
        <Card className="m-3">
            <CardContent>
                <AspectRatio ratio={16 / 9} className="w-3/4">
                    <video
                    className="w-full h-full rounded-md object-cover"
                    controls
                    preload="metadata"
                    >  Your browser does not support the video tag.</video>
                </AspectRatio>
            </CardContent>
            <CardContent>
                <CardTitle>
                    Loading....
                </CardTitle>
            </CardContent>
        </Card>
    )
}