import {
  Card,
  CardDescription,
  CardContent,
  CardTitle,
} from "@/components/ui/card"

import {Badge} from "@/components/ui/badge"
import { AspectRatio } from "@/components/ui/aspect-ratio"
import { Skeleton } from "@/components/ui/skeleton"

type VideoData = {
    row_id: string,
    s3_id: string,

    custom_title: string | null,
    custom_description: string | null,

    username: string,
    filename: string,
    video_url: string,
}

const DESCRIPTION_MAX_CHAR = 300

function limitString(text: string): string {
  return text.length > DESCRIPTION_MAX_CHAR ? text.slice(0, DESCRIPTION_MAX_CHAR) : text;
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
                <CardDescription> 
                    <Badge variant="secondary">{props.username}</Badge> 
                    {props.custom_description ? limitString(props.custom_description) : "No description has been given"}
                </CardDescription>
                
            </CardContent>
        </Card>
    )
}

export const SkeletonVideoCard = () => {
    return (
        <Card className="m-3">
            <CardContent>
                <AspectRatio ratio={16 / 9} className="w-3/4">
                    <Skeleton className="h-full w-full" />
                </AspectRatio>
            </CardContent>
            <CardContent>
                <Skeleton className="h-5 w-[200px]" />
                <Skeleton className="h-4 w-[400px]" />
            </CardContent>
        </Card>
    )
}