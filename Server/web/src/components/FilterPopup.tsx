
import { useState, useEffect } from 'react';

import {
    useQuery,
} from '@tanstack/react-query'

import {
    Command,
    CommandEmpty,
    CommandGroup,
    CommandInput,
    CommandItem,
    CommandList,
} from "@/components/ui/command"

import {
  Field,
  FieldGroup,
} from "@/components/ui/field"

import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogHeader,
    DialogTitle,
} from "@/components/ui/dialog"

import {
    Button
} from "@/components/ui/button"

import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"

import {Badge} from "@/components/ui/badge"

import { useFilter } from '@/context/FilterContext';
import { Check } from 'lucide-react';

const TITLE_DEBOUNCE = 500
type Tag = {
    tag_id: string,
    name: string,
    color: string,
}

type TagProps = {
    tags: Array<Tag>
    open: boolean
    setOpen: (open: boolean) => void
}

const FETCH_TAGS_ENDPOINT = "/api/v2/tags/get";

function snakeToTitleCase(str: string) {
  return str
    .split('_')
    .map(word => word.charAt(0).toUpperCase() + word.slice(1).toLowerCase())
    .join(' ');
}

export const FilterPopup = (props: TagProps) => {
    const {filter, toggleTag, updateFilterParameter, cycleFilterDirection, cycleFilterElement} = useFilter()

    const [title, setTitle] = useState("")
    const [debouncedTitle, setDebouncedTitle] = useState("");

    useEffect(() => {
        const timer = setTimeout(() => {
            setDebouncedTitle(title);
        }, TITLE_DEBOUNCE); 

        return () => clearTimeout(timer); 
    }, [title]);

    useEffect(() => {
        updateFilterParameter("title", debouncedTitle)
    }, [debouncedTitle]);

    const { data } = useQuery<Tag[]>({
        queryKey: [`get-all-tags`],
        queryFn: () =>
        fetch(`${FETCH_TAGS_ENDPOINT}`).then((res) => {
            if (!res.ok) throw new Error(`Failed to fetch tags: ${res.status}`)
            return res.json()
        }),
    })

    return (
        <>
        <Dialog open={props.open} onOpenChange={props.setOpen}>
            <DialogContent className="overflow-hidden">
                <DialogHeader>
                    <DialogTitle>Change Filter Settings</DialogTitle>
                        <DialogDescription>
                        Modify the filter for videos.
                        </DialogDescription>
                </DialogHeader>
                <FieldGroup>
                    <Field>
                        <Label htmlFor="title-1">Video Title</Label>
                        <Input id="title-1" onChange={(e) => setTitle(e.target.value)}name="title" placeholder="" defaultValue="" />
                    </Field>
                    <Field>
                        <Label htmlFor="element-1">Filter By</Label>
                        <Button variant="outline" onClick={() => {cycleFilterElement(filter.filter_element)}}> {snakeToTitleCase(filter.filter_element)} </Button>
                    </Field>
                    <Field>
                        <Label htmlFor="direction-1">Ordering Direction</Label>
                        <Button variant="outline"  onClick={() => {cycleFilterDirection(filter.filter_direction)}}> {snakeToTitleCase(filter.filter_direction)} </Button>
                    </Field>
                </FieldGroup>
                <Command>
                    <CommandInput placeholder="Search for a tag..." />
                    <CommandList>
                        <CommandEmpty>No results found.</CommandEmpty>
                        <CommandGroup>
                        
                        {data?.map((tag: Tag) => (
                            <CommandItem value={tag.name} onSelect={() => {
                                toggleTag(tag.tag_id)
                            }}>
                                <div className="flex flex-row items-center"> 
                                    <Badge variant="secondary" style={{backgroundColor: tag.color}}>{tag.name}</Badge> 
                                    { filter.filter_tags.has(tag.tag_id) ? <Check className="ml-1"/> : null}
                                </div>
                            </CommandItem>
                        ))}

                        </CommandGroup>
                    </CommandList>
                </Command>
            </DialogContent>

            {/* <div className="border-t p-2">
                <Button className="w-full" onClick={openCreateMenu}> <Plus className="mr-2 h-4 w-4" /> Create Tag</Button>
            </div> */}
        </Dialog>
        {/* <Dialog open={createOpen} onOpenChange={setCreateOpen}>
            <DialogContent className="sm:max-w-sm">        
                <form onSubmit={handleSubmit}>
                    <DialogHeader>
                        <DialogTitle>Create Tag</DialogTitle>
                        <DialogDescription>
                            Add a new global tag.
                        </DialogDescription>
                    </DialogHeader>
                    <FieldGroup className='mt-1'>
                        <Field>
                        <Input className={validInput ? "mb-3" : "mt-1"} id="tag-1" name="tag" placeholder="New Tag..." required aria-invalid={!validInput}/>                        
                        {!validInput && (
                            <FieldDescription className='mb-3'>
                                This tag already exists.
                            </FieldDescription>                            
                        )}
                        </Field>
                    </FieldGroup>
                    <DialogFooter>
                        <Button type="button" variant="outline" onClick={closeCreateMenu}>Cancel</Button>
                        <Button type="submit">Save</Button>
                    </DialogFooter>        
                </form>
            </DialogContent>
        </Dialog> */}
        </>
    )
}