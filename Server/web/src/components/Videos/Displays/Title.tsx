import { Button } from "@/components/ui/button"
import { Pencil, Check } from "lucide-react"
import { Input } from "@/components/ui/input"
import { Field } from "@/components/ui/field"

import {
  useMutation,
  useQueryClient
} from '@tanstack/react-query'

import {
  CardTitle,
} from "@/components/ui/card"
import { useState } from "react"
import { useAuth } from "@/context/AuthContext"

type VideoDataProps = {
    title?: string | null,
    filename: string,
    id: string,
    can_modify: boolean
}

const ENDPOINT = "/api/v2/videos/title";

const MAX_TITLE_LENGTH = 60

function limitStringLength(input: string) {
    console.log(input)
    return input.length > MAX_TITLE_LENGTH ? input.substring(0, MAX_TITLE_LENGTH) : input 
}

export const VideoTitleDisplay = (props: VideoDataProps) => {
    const { session, authLoading } = useAuth()
    const [isUpdating, setUpdating] = useState(false)
    const [title, setTitle] = useState(props.title ? limitStringLength(props.title) : limitStringLength(props.filename))
    const queryClient = useQueryClient()
    
    const updateTitleRequest = useMutation<null, Error, {row_id: string, title: string}>({
		mutationKey: [`update-video-title`, props.id],
		mutationFn: (payload: any) =>
        fetch(`${ENDPOINT}`, {
            method: "PUT",
            body: JSON.stringify(payload),
            headers: {
                "Content-Type": "application/json",
                "Authorization": `Bearer ${session?.access_token}`
            },
        }).then((res) => {
            if (!res.ok) throw new Error(`Failed to fetch tags: ${res.status}`)
            return res.json()
        }),
		onError: (err) => {
			console.log(err)			
		},
        onSuccess: () => {
            queryClient.invalidateQueries({
                queryKey: [`get-video-data`, props.id]
            })
        },
		onSettled: () => {
			setUpdating(false)
		}
	})
    function startUpdating() {
        setUpdating(true)
    }

    function stopUpdating() {
        if (!title || title.length == 0){
            setTitle(limitStringLength(props.filename))
            return
        }

        if (authLoading) {
            return //TODO: Send Error Popup
        }
        const body = {
            row_id: props.id,
            title: title
        }
        updateTitleRequest.mutate(body)
    }

    function updateTitle(newTitle: string) {
        setTitle(limitStringLength(newTitle))
    }

    return (
        <div className="w-full flex-1">
            <Field orientation="horizontal">
                {(!isUpdating && 
                    <div className="flex flex-row w-full">
                        <CardTitle className="line-clamp-1 truncate items-center">
                            {title}
                        </CardTitle>
                        {(props.can_modify) &&
                            <Button className="shrink-0" variant="outline" size="icon-xs" aria-label="Change Title" onClick={startUpdating}><Pencil /></Button>   
                        }        
                    </div>
                )}
                {isUpdating &&
                    <Field className="flex items-center mb-2" orientation="horizontal">
                        <Input 
                            placeholder="Enter title..."
                            value={title}
                            className="max-w-[400px] w-3/4"
                            maxLength={MAX_TITLE_LENGTH}                       
                            onChange={(e) => updateTitle(e.target.value)} 
                            onKeyDown={(e) => {
                                if (e.key === "Enter") {
                                    stopUpdating();
                                }
                            }} /> 
                        <Button variant="outline" size="icon-sm" aria-label="Confirm Title" onClick={stopUpdating}><Check /></Button>
                    </Field>
                }
            </Field>
        </div>
    )
}