import {
  Card,
  CardContent,
} from "@/components/ui/card"

import { User, Eye, EyeOff } from "lucide-react"

import { ScrollArea} from "@/components/ui/scroll-area"
import {Badge} from "@/components/ui/badge"
import { AspectRatio } from "@/components/ui/aspect-ratio"
import { Skeleton } from "@/components/ui/skeleton"

import { VideoTitleDisplay } from "./Displays/Title"
import { VideoDescDisplay } from "./Displays/Description"

import { useMediaQuery } from "@/hooks/MediaQuery"

import { memo, useState } from "react"
import { TagDisplay } from "./Displays/Tags"
import { VideoScrollBar } from "./ScrollBar/VideoScrollBar"
import { Button } from "@base-ui/react"
import { ChangeVisibilityPopup } from "./ChangeVisibility"
 
type VideoData = {
    row_id: string,
    s3_id: string,

    custom_title?: string | null,
    custom_description?: string | null,

    visibility: string,

    username: string,
    filename: string,

    video_url: string,
    thumbnail_url: string,

    timestamp: number,
    can_modify: boolean
}



const VideoPlayer = memo(({ src, thumbnail }: { src: string, thumbnail: string }) => (
    
    <video
        className="w-full h-full rounded-md object-cover"
        src={src}
        controls
        autoPlay
        poster={thumbnail}
        preload="auto"
    > 
        Your browser does not support the video tag.
    </video>
));

export const VideoCard = (props: VideoData) => {
    const [visibility, setVisibility] = useState<string>(props.visibility)
    const [visibilityMenuOpen, setVisibilityMenuOpen] = useState<boolean>(false)
    const isVertical = useMediaQuery("(max-width: 903px)")
    console.log(isVertical)

    return (
        <div className="flex w-full h-full">
            <ChangeVisibilityPopup key={props.row_id} video_id={props.row_id} visibility={props.visibility} open={visibilityMenuOpen} setOpen={setVisibilityMenuOpen} onSelect={(newVisibility: string) => {setVisibility(newVisibility)}}/>
            <Card className={isVertical ? "bg-gray w-full flex flex-col gap-0 mx-3" : "bg-gray w-full gap-0 flex flex-row mx-3"}>
                <div className={isVertical ? "w-full" : "h-full w-3/4"}>
                    <CardContent>
                        <AspectRatio ratio={16 / 9} className="w-full">
                            <VideoPlayer src={props.video_url} thumbnail={props.thumbnail_url}/> 
                        </AspectRatio>
                    </CardContent>
                    <CardContent className={isVertical ? "pt-3" : "pt-3 h-full"}>
                        <div className="flex flex-row justify-between items-center">
                            <VideoTitleDisplay key={props.row_id} can_modify={props.can_modify} title={props.custom_title} filename={props.filename} id={props.row_id}/>
                        </div>                            
                        <VideoDescDisplay key={props.row_id} can_modify={props.can_modify} description={props.custom_description} id={props.row_id}/>
                        <div className="shrink-0 flex mr-1">
                            <Badge variant="secondary" className=""><User/>{props.username}</Badge>
                            {(visibility && 
                            (
                                <>
                                {(props.can_modify) ? (
                                    <Button onClick={() => {setVisibilityMenuOpen(true)}}><Badge variant="secondary" className="ml-2"> {(visibility== "Public" ? (<Eye/>) : (<EyeOff/>))}{visibility}</Badge></Button>
                                ) : (
                                    <Badge variant="secondary" className="ml-2"> {(visibility== "Public" ? (<Eye/>) : (<EyeOff/>))}{visibility}</Badge>
                                )}
                                </>
                            ))}                           
                        </div>
                        <TagDisplay video_id={props.row_id} can_modify={props.can_modify}/>
                    </CardContent>
                </div>
                {isVertical ? ( 
                    <div className="w-full h-full min-h-0 px-2 mt-2">
                        <ScrollArea className="h-full w-full min-h-0 rounded-md border">
                            <VideoScrollBar key={props.row_id} cursor={{
                                "row_id": props.row_id, 
                                "timestamp":props.timestamp
                            }}/>
                        </ScrollArea>
                    </div>
                ) : (
                    <div className="w-1/4 mr-3 h-full">
                        <ScrollArea className="h-full w-full min-h-0 rounded-md border">
                            <VideoScrollBar key={props.row_id} cursor={{
                                "row_id": props.row_id, 
                                "timestamp":props.timestamp
                            }}/>
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