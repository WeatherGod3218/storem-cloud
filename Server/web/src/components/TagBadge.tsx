import {Badge} from "@/components/ui/badge"
import { X } from "lucide-react"

type Tag = {
    tag_id: string,
    name: string,
    color: string,
}

type TagProps = {
    tag: Tag,
    onRemove: (tag: Tag) => void    
    cantRemove?: boolean
}


export const TagBadge = (props: TagProps) => {
    return (
        <div className=" mr-1 group relative inline-flex items-center">
            <Badge variant="secondary" style={{ backgroundColor: props.tag.color}}>
                {props.tag.name}
                {!props.cantRemove && (
                    <button
                        onClick={() => props.onRemove(props.tag)}
                        aria-label={`Remove ${props.tag.name}`}
                        className="opacity-0 group-hover:opacity-100 focus-visible:opacity-100 transition-opacity"
                    >
                        <X className="h-3 w-3" />
                    </button>
                )}
            </Badge>
            {/* <button
                onClick={() => onRemove(tag)}
                aria-label={`Remove ${tag.name}`}
                className="ml-1 opacity-0 group-hover:opacity-100 focus-visible:opacity-100 transition-opacity"
            >
                <X className="h-3 w-3" />
            </button> */}
        </div>
    )
}