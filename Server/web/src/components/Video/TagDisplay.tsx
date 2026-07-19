import {
  useQuery,
  useMutation,
  useQueryClient
} from '@tanstack/react-query'

import { TagPopup } from "../TagPopup"
import { TagBadge } from "../TagBadge"
import { useState} from "react"

import { Button } from "@/components/ui/button"
import { Plus } from "lucide-react"

type TagProps = {
    video_id: string
}

type Tag = {
    tag_id: string
    name: string
    color: string
}

const ADD_TAG_ENDPOINT = "/api/v2/tags/video/add";
const FETCH_TAGS_ENDPOINT = "/api/v2/tags/video/get";
const REMOVE_TAG_ENDPOINT = "/api/v2/tags/video/remove";


export const TagDisplay = (props: TagProps) => {
    const [open, setOpen] = useState(false)
    const queryClient = useQueryClient()

    const removeMutation = useMutation({
        mutationKey: [`delete-video-tag`, props.video_id],
        mutationFn: (tag: Tag) =>
        fetch(`${REMOVE_TAG_ENDPOINT}`, {
            method: "DELETE",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({
                video_id: props.video_id,
                tag_id: tag.tag_id
                
            })
        }).then((res) => {
            if (!res.ok) {
                throw new Error(`Request failed: ${res.status}`)
            }
            res.json()
        }),
        onError: (err) => {
            console.error(err)
        },
        onSuccess: () => {
            queryClient.invalidateQueries({
                queryKey: [`get-video-tags`, props.video_id]
            })
        }
    })

    const addMutation = useMutation({
        mutationKey: [`add-video-tag`, props.video_id],
        mutationFn: (tag: Tag) =>
        fetch(`${ADD_TAG_ENDPOINT}`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({
                video_id: props.video_id,
                tag_id: tag.tag_id
            })
        }).then((res) => {
            if (!res.ok) {
                throw new Error(`Request failed: ${res.status}`)
            }
            res.json()
        }),
        onError: (err) => {
            console.error(err)
        },
        onSuccess: () => {
            queryClient.invalidateQueries({
                queryKey: [`get-video-tags`, props.video_id]
            })
        }
    })

    const { isPending, error, data } = useQuery<Tag[]>({
        queryKey: [`get-video-tags`, props.video_id],
        queryFn: () =>
        fetch(`${FETCH_TAGS_ENDPOINT}/${props.video_id}`).then((res) => {
            if (!res.ok) throw new Error(`Failed to fetch tags: ${res.status}`)
            return res.json()
        }),
    })

    console.log(data)

    const tags = new Map(data?.map(tag => [tag.tag_id, tag]))
    
    function addTagToVideo(tag: Tag) {
        console.log(tag)
        addMutation.mutate(tag)
    }

    return (       
        <> 
        <TagPopup open={open} setOpen={setOpen} onSelect={addTagToVideo} tags={[...tags.values()]}/>
        <div className="pb-2 pt-1 flex flex-row items-center">
            {isPending ? (
                <h6 className="flex text-muted-foreground">Loading Tags...</h6>   
            ) : error ? (
                <h4 className="flex text-destructive">Failed To Load Tags!</h4>   
            ) : (
            <>

            {tags.size == 0 ? (
                <h6 className="flex text-muted-foreground">No tags on this video.</h6>   
            ) : (
                [...tags.values()].map((tag: Tag) => (
                    <TagBadge key={tag.tag_id} tag={tag} onRemove={(tag) => {
                        removeMutation.mutate(tag)
                    }}/>
                ))
            )}

            <Button 
                className="ml-1 text-muted-foreground" 
                variant="outline" 
                onClick={() => {setOpen(true)}} 
                size="icon-xs" 
                aria-label="Add Tag"
            >
                <Plus />
            </Button>
            </>
            )}
        </div>
        </>    
    )
}