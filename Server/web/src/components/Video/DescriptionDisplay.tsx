import { Button } from "@/components/ui/button"
import { Pencil, Check } from "lucide-react"
import { Input } from "@/components/ui/input"
import { Field } from "@/components/ui/field"
import { CardDescription } from "@/components/ui/card"

import { useState } from "react"

type VideoDataProps = {
    description: string | null,
    id: string,
}

const DEFAULT_DESCRIPTION = "No description has been given."
const ENDPOINT = "/api/v2/videos/description";

export const VideoDescDisplay = (props: VideoDataProps) => {
    const [isUpdating, setUpdating] = useState(false)
    const [desc, setDesc] = useState(props.description ? props.description : DEFAULT_DESCRIPTION)
    // const [description, setDescription] = useState(tempDesc)
    
    function startUpdating() {
        setUpdating(true)
    }

    function stopUpdating() {
        let cancelled = false
        const body = {
            row_id: props.id,
            description: desc
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

    function updateTitle(newDesc: string) {
        if (!newDesc){
            setDesc(DEFAULT_DESCRIPTION)
            return
        }
        setDesc(newDesc)
    }

    return (
        <div>
            <Field orientation="horizontal">
                {(!isUpdating && 
                    <CardDescription className="justify-center items-center">
                        {desc}   <Button variant="outline" size="icon-xs" aria-label="Change Description" onClick={startUpdating}><Pencil /></Button>
                    </CardDescription>
                )}
                {isUpdating &&
                    <Field className="flex items-center mb-2" orientation="horizontal">
                        <Input 
                            placeholder="Enter description..." 
                            value={desc}
                            className="max-w-[400px] w-3/4"
                            maxLength={200}                                                
                            onChange={(e) => updateTitle(e.target.value)} 
                            onKeyDown={(e) => {
                                if (e.key === "Enter") {
                                    stopUpdating();
                                }
                            }} /> 
                        <Button variant="outline" size="icon-sm" aria-label="Confirm Description" onClick={stopUpdating}><Check /></Button>
                    </Field>
                }
            </Field>
        </div>
    )
}