import {
  Card,
  CardContent,
} from "@/components/ui/card"

import { ScrollArea} from "@/components/ui/scroll-area"
import {Badge} from "@/components/ui/badge"
import { AspectRatio } from "@/components/ui/aspect-ratio"
import { Skeleton } from "@/components/ui/skeleton"

import { VideoTitleDisplay } from "./TitleDisplay"
import { VideoDescDisplay } from "./DescriptionDisplay"

import { useMediaQuery } from "@/hooks/MediaQuery"

import { memo } from "react"
import { TagDisplay } from "./TagDisplay"
 
type VideoData = {
    row_id: string,
    s3_id: string,

    custom_title: string | null,
    custom_description: string | null,

    username: string,
    filename: string,
    video_url: string,
}

// const DESCRIPTION_MAX_CHAR = 300

// function limitString(text: string): string {
//   return text.length > DESCRIPTION_MAX_CHAR ? text.slice(0, DESCRIPTION_MAX_CHAR) : text;
// }

const VideoPlayer = memo(({ src }: { src: string }) => (
    
    <video
        className="w-full h-full rounded-md object-cover"
        src={src}
        controls
        preload="metadata"
    > 
        Your browser does not support the video tag.
    </video>
));

const videos = Array.from({ length: 50 }).map(
  (_, i, __) => `video number ${i}`
)

export const VideoCard = (props: VideoData) => {
    const isVertical = useMediaQuery("(max-width: 903px)")
    console.log(isVertical)

    return (
        <div className="flex w-full h-full">
            <Card className={isVertical ? "bg-gray w-full flex flex-col gap-0 mx-3" : "bg-gray w-full gap-0 flex flex-row mx-3"}>
                <div className={isVertical ? "w-full" : "h-full w-3/4"}>
                    <CardContent>
                        <AspectRatio ratio={16 / 9} className="w-full">
                            <VideoPlayer src={props.video_url}/> 
                        </AspectRatio>
                    </CardContent>
                    <CardContent className={isVertical ? "pt-3" : "pt-3 h-full"}>
                        <div className="flex justify-between items-center">
                            <VideoTitleDisplay key={props.row_id} title={props.custom_title} filename={props.filename} id={props.row_id}/>
                            <Badge variant="secondary" className="mr-1">{props.username}</Badge> 
                        </div>                            
                        <VideoDescDisplay key={props.row_id} description={props.custom_description} id={props.row_id}/>
                        <TagDisplay video_id={props.row_id}/>              
                    </CardContent>
                </div>
                {isVertical ? ( 
                    <div className="w-full h-full min-h-0 px-2 mt-2">
                        <ScrollArea className="h-full w-full min-h-0 rounded-md border">
                            <div className="p-4">
                                {videos.map((vidname) => (
                                    <h4 key={vidname} className="mb-4 text-sm leading-none font-medium">{vidname}</h4>   
                                ))}
                            </div>
                        </ScrollArea>
                    </div>
                ) : (
                    <div className="w-1/4 mr-3 h-full">
                        <ScrollArea className="h-full w-full min-h-0 rounded-md border">
                            <div className="p-4">
                                {videos.map((vidname) => (
                                    <h4 key={vidname} className="mb-4 text-sm leading-none font-medium">{vidname}</h4>   
                                ))}
                            </div>
                        </ScrollArea>
                    </div>
                )}
            </Card>              
        </div>
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