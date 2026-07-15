import { useState } from "react"

import {
    Card,
    CardAction,   
    CardDescription,
    CardFooter,
    CardHeader,
    CardTitle,
} from "@/components/ui/card"
import {Badge} from "@/components/ui/badge"
import { Skeleton } from "@/components/ui/skeleton"

import { useNavigate } from "react-router"
import { cn } from "@/lib/utils"
// type props = {
//     test: string | null,
// }

type ThumbnailCardProps = {
    rowId: string,

    username: string,

    filename: string,
    thumbnail: string,

    customTitle?: string | null,
    customDescription?: string | null,
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
    
    console.log(props.thumbnail)

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
            src={props.thumbnail}
            alt="Video Cover"
            onLoad={() => setImageLoaded(true)}
            className={cn(
                "relative z-20 aspect-video w-full object-cover brightness-100 dark:brightness-100",
                imageLoaded ? "opacity-100" : "opacity-0"
            )}
            />
        <CardHeader>
            <CardTitle className="w-2/3">{props.customTitle ? props.customTitle : props.filename}</CardTitle>
            <CardAction>
                
            </CardAction>
            <CardDescription>{props.customDescription ? limitString(props.customDescription) : "No description has been given"}</CardDescription>
            <Badge variant="secondary">{props.username}</Badge>
        </CardHeader>
        <CardFooter>
            <Badge className="bg-blue-50 text-blue-700 dark:bg-blue-950 dark:text-blue-300">
                Tag1
            </Badge>
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