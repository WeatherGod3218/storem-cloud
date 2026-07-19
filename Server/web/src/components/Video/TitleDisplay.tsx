import { Button } from "@/components/ui/button"
import { Pencil, Check } from "lucide-react"
import { Input } from "@/components/ui/input"
import { Field } from "@/components/ui/field"

import {
  CardTitle,
} from "@/components/ui/card"
import { useState } from "react"

type VideoDataProps = {
    title: string | null,
    filename: string,
    id: string,
}

const ENDPOINT = "/api/v2/videos/title";

const MAX_TITLE_LENGTH = 45

function limitStringLength(input: string) {
    console.log(input)
    return input.length > MAX_TITLE_LENGTH ? input.substring(0, MAX_TITLE_LENGTH) : input 
}

export const VideoTitleDisplay = (props: VideoDataProps) => {
    const [isUpdating, setUpdating] = useState(false)
    const [title, setTitle] = useState(props.title ? limitStringLength(props.title) : limitStringLength(props.filename))
    
    function startUpdating() {
        setUpdating(true)
    }

    function stopUpdating() {
        let cancelled = false
        const body = {
            row_id: props.id,
            title: title
        }

        fetch(`${ENDPOINT}`, {
            method: "PUT",
            body: JSON.stringify(body),
            headers: {
                "Content-Type": "application/json",
            },
        })
        .then(res => {if (!res.ok) {throw new Error(`HTTP ERROR: ${res.status}`)} return res.json()})
        .catch(err => { if (!cancelled) {console.log(`${err}`)}})
        .finally(() => { if (!cancelled) setUpdating(false); })
        return () => { cancelled = true; };
    }

    function updateTitle(newTitle: string) {
        if (!newTitle){
            setTitle(limitStringLength(props.filename))
            return
        }
        setTitle(limitStringLength(newTitle))
    }

    return (
        <div className="w-4/5 ">
            <Field orientation="horizontal">
                {(!isUpdating && 
                    <CardTitle className="justify-center items-center">
                        {title}   <Button variant="outline" size="icon-xs" aria-label="Change Title" onClick={startUpdating}><Pencil /></Button>
                    </CardTitle>
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