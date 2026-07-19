import { Button } from "@/components/ui/button"
import { Pencil, Check } from "lucide-react"
import { Input } from "@/components/ui/input"
import { Field } from "@/components/ui/field"

import {
  CardTitle,
} from "@/components/ui/card"
import { useState } from "react"

const tempTitle = "Test Title"
// const tempDesc = "Test Desc"

export const TestPage = () => {
    const [isUpdating, setUpdating] = useState(true)
    const [title, setTitle] = useState(tempTitle)
    // const [description, setDescription] = useState(tempDesc)
    
    function startUpdating() {
        setUpdating(true)
    }

    function stopUpdating() {
        console.log("MEOW!!!")
        setUpdating(false)
    }

    return (
        <div className="w-full min-h-screen flex justify-center items-center">
            <Field orientation="horizontal">
                {(!isUpdating && 
                    <CardTitle className="justify-center items-center">
                        {title}   <Button variant="outline" size="icon-xs" aria-label="Change Title" onClick={startUpdating}><Pencil /></Button>
                    </CardTitle>
                )}
                {isUpdating &&
                    <Field orientation="horizontal">
                        <Input placeholder="Enter title..." onChange={(e) => setTitle(e.target.value)} className="w-full"/> 
                        <Button variant="outline" size="icon-sm" aria-label="Confirm Title" onClick={stopUpdating}><Check /></Button>
                    </Field>
                }
            </Field>
        </div>
    )
}