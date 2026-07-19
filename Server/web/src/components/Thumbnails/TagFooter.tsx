import {
  useQuery,
} from '@tanstack/react-query'

// import { Spinner } from "@/components/ui/spinner"
import { TagBadge } from "../TagBadge"

type TagProps = {
    video_id: string
}

type Tag = {
    tag_id: string
    name: string
    color: string
}

const FETCH_TAGS_ENDPOINT = "/api/v2/tags/video/get";

export const TagFooter = (props: TagProps) => {
    const { isPending, error, data } = useQuery<Tag[]>({
        queryKey: [`get-video-tags-thumbnails`, props.video_id],
        queryFn: () =>
        fetch(`${FETCH_TAGS_ENDPOINT}/${props.video_id}`).then((res) => {
            if (!res.ok) throw new Error(`Failed to fetch tags: ${res.status}`)
            return res.json()
        }),
    })

    console.log(data)


    return (       
        <> 
        <div className="pb-2 pt-1 flex flex-row items-center">
            {isPending ? (
                <h6 className="flex text-muted-foreground">Loading Tags...</h6>   
            ) : error ? (
                <h4 className="flex text-destructive">Failed To Load Tags!</h4>   
            ) : (
            <>
                {data.length == 0 ? (
                    <h6 className="flex text-muted-foreground">No tags.</h6>   
                ) : (
                    data.map((tag: Tag) => (
                    // data.slice(0, 4).map((tag: Tag) => (
                        <TagBadge key={tag.tag_id} tag={tag} cantRemove onRemove={(_) => {}}/>
                    ))
                )}
            </>
            )}
        </div>
        </>    
    )
}