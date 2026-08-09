import { useState } from "react"
import { useNavigate } from "react-router"

import { Skeleton } from "@/components/ui/skeleton"

import {
  Item,
  ItemContent,
  ItemDescription,
  ItemGroup,
  ItemMedia,
  ItemTitle,
} from "@/components/ui/item"

import { cn } from "@/lib/utils"


type ScrollBarThumbnailProps = {
    rowId: string,

    username: string,

    filename: string,
    thumbnail: string,

    customTitle?: string | null,
    customDescription?: string | null,
}

export function ScrollBarThumbnail(props: ScrollBarThumbnailProps) {
    const [imageLoaded, setImageLoaded] = useState(false);
	const navigate = useNavigate()

	const handleSelect = () => {
        navigate(`/video/${props.rowId}`)
    }

	return (
	<div className="pt-3 flex w-full flex-col">
		<ItemGroup className="gap-2">
			<Item
			key={props.rowId}
			variant="outline"
			role="listitem"
			onClick={handleSelect}
			render={
				<a href="#" className="flex items-center gap-3">
					
				<ItemMedia variant="image" className="shrink-0">
					{!imageLoaded && (
						<Skeleton className="object-cover absolute inset-0 h-[32px] w-[40px] h-full"/>
					)}
					<img
					src={props.thumbnail}
					alt="Video Thumbnail"
					width={200}
					height={200}
					onLoad={() => setImageLoaded(true)}
					className={cn(
						"object-cover",
						imageLoaded ? "opacity-100" : "opacity-0"
					)}
					/>
				</ItemMedia>
				<ItemContent className="flex flex-row min-w-0 items-center justify-between">
					<ItemTitle className="min-w-0 flex line-clamp-2 break-words">
					{props.customTitle ? props.customTitle : props.filename}
					</ItemTitle>
					<ItemDescription className="shrink-0 ext-right">{props.username}</ItemDescription>
				</ItemContent>
				</a>
			}
			/>
		</ItemGroup>
    </div>
  	)
}
