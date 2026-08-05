import { useState } from "react"
import { useNavigate } from "react-router"

import {
    Card,
    CardAction,   
    CardDescription,
    CardFooter,
    CardHeader,
    CardTitle,
} from "@/components/ui/card"

import { Eye, EyeOff, User } from "lucide-react"

import {Badge} from "@/components/ui/badge"
import { Skeleton } from "@/components/ui/skeleton"

import { cn } from "@/lib/utils"
import { TagFooter } from "./TagFooter"

type ThumbnailCardProps = {
    rowId: string,

    username: string,

    filename: string,
    thumbnail: string,

    customTitle?: string | null,
    customDescription?: string | null,

    visibility?: string
}

const DESCRIPTION_MAX_CHAR = 100

function limitString(text: string): string {
  return text.length > DESCRIPTION_MAX_CHAR ? text.slice(0, DESCRIPTION_MAX_CHAR) : text;
}

export const ThumbnailCard = (props: ThumbnailCardProps) => {
    const [imageLoaded, setImageLoaded] = useState(false);
    const navigate = useNavigate()

    const handleSelect = () => {
        navigate(`/video/${props.rowId}`)
    }

    return (
    <Card 
        onClick={handleSelect} 
        role="button"
        tabIndex={0}
        className={cn(
                "relative mx-auto w-full max-w-sm pt-0 cursor-pointer",
                "transition-colors hover:border-primary/50 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary"
            )}
        >        
        <div className="absolute inset-0 z-30 aspect-video" />
            {!imageLoaded && (
                <Skeleton className="absolute inset-0 w-full h-full" />
            )}
            <img
            src={props.thumbnail || undefined}
            alt="Video Cover"
            onLoad={() => setImageLoaded(true)}
            className={cn(
                "relative z-20 aspect-video w-full object-cover brightness-100 dark:brightness-100",
                imageLoaded ? "opacity-100" : "opacity-0"
            )}
            />
        <CardHeader>
            <CardTitle className="w-full min-w-0 line-clamp-2 break-words">{props.customTitle ? props.customTitle : props.filename}</CardTitle>
            <CardAction>
            </CardAction>
            <CardDescription>{props.customDescription ? limitString(props.customDescription) : "No description has been given"}</CardDescription>
            <Badge variant="secondary" className=""><User/>{props.username}</Badge>
            {(props.visibility && 
            (
                props.visibility == "Public" ? (
                < Badge variant="secondary" className="ml-2"> <Eye/>{props.visibility}</Badge>
            ) : (
                <Badge variant="secondary" className="ml-2">  <EyeOff/>{props.visibility}</Badge>
            )))} 
        </CardHeader>
        <div className="flex-1"/> 
        <CardFooter>
            <TagFooter key={props.rowId} video_id={props.rowId}/>
        </CardFooter>
    </Card>
    )
}

export const ThumbnailSkeletonCard = () => {
    const navigate = useNavigate()

    const handleSelect = () => {
        navigate("/video")
    }
    
    return (
    <Card 
        onClick={handleSelect} 
        role="button"
        tabIndex={0}
        className={cn(
                "relative mx-auto w-full max-w-sm pt-0 cursor-pointer",
                "transition-colors hover:border-primary/50 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary"
            )}
        >        
        <div className="absolute inset-0 z-30 aspect-video" />
            <img 
            src=""
            alt=""
            className="relative z-20 aspect-video w-full object-cover brightness-100 dark:brightness-100 dark:bg-zinc-800"
            />
        <CardHeader className="gap-2">
            <Skeleton className="h-5 w-[200px]" />
            <Skeleton className="h-4 w-[270px]" />
        </CardHeader>
        <CardFooter className="gap-1">
            <Skeleton className="h-5 w-[50px]" />
            <Skeleton className="h-5 w-[50px]" />
            <Skeleton className="h-5 w-[50px]" />
        </CardFooter>
    </Card>
    )
}