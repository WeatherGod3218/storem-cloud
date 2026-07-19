
import { useState, useCallback } from 'react';

import {
    useQuery,
    useMutation,
    useQueryClient
} from '@tanstack/react-query'

import {
    Command,
    CommandDialog,
    CommandEmpty,
    CommandGroup,
    CommandInput,
    CommandItem,
    CommandList,
} from "@/components/ui/command"

import {
  Field,
  FieldGroup,
  FieldDescription
} from "@/components/ui/field"

import {
    Dialog,
    DialogFooter,
    DialogContent,
    DialogDescription,
    DialogHeader,
    DialogTitle,
} from "@/components/ui/dialog"

import { Input } from "@/components/ui/input"

import { Plus } from "lucide-react"
import { Button } from "@/components/ui/button"
import {Badge} from "@/components/ui/badge"

import { Fragment } from 'react';

type Tag = {
    tag_id: string,
    name: string,
    color: string,
}

type TagProps = {
    tags: Array<Tag>
    open: boolean
    setOpen: (open: boolean) => void
    onSelect: (tag: Tag) => void
}

const FETCH_TAGS_ENDPOINT = "/api/v2/tags/get";
const CREATE_TAG_ENDPOINT = "/api/v2/tags/create";

function useTagSet(tags: Tag[]) {
  const existing = new Set(tags.map(t => t.tag_id));
  const has = useCallback((tag: Tag) => existing.has(tag.tag_id), [tags]);
  return { has };
}


export const TagPopup = (props: TagProps) => {
    const [validInput, setValidInput] = useState<boolean>(true)
    const {has} = useTagSet(props.tags)
    const [createOpen, setCreateOpen] = useState<boolean>(false)
    const queryClient = useQueryClient()

    const createMutation = useMutation({
        mutationKey: [`create-tag`],
        mutationFn: (tagName: string) =>
        fetch(`${CREATE_TAG_ENDPOINT}`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({
                name: tagName,
            })
        }).then((res) => {
            if (!res.ok) {
                if (res.status == 409) {
                    setValidInput(false)
                }
            }
            closeCreateMenu()
        }),
        onSuccess: () => {
            queryClient.invalidateQueries({
                queryKey: [`get-all-tags`]
            })
        }    
    })

    function handleSubmit(e: React.FormEvent<HTMLFormElement>) {
        e.preventDefault();
        const formData = new FormData(e.currentTarget);
        const tagName = formData.get("tag") as string;

        if (!tagName.trim()) {
            setValidInput(false) 
            return
        }
    
        createMutation.mutate(tagName);
    }

    function openCreateMenu() {
        props.setOpen(false)
        setCreateOpen(true)
    }

    function closeCreateMenu() {
        props.setOpen(true)
        setCreateOpen(false)
    }

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
        <CommandDialog open={props.open} onOpenChange={props.setOpen}>
            <Command>
                <CommandInput placeholder="Search for a tag..." />
                <CommandList>
                    <CommandEmpty>No results found.</CommandEmpty>
                    <CommandGroup>
                    
                    {data?.map((tag: Tag) => (
                        <Fragment key={tag.tag_id}>
                        {!has(tag) && (
                            <CommandItem value={tag.tag_id} onSelect={() => {
                                props.onSelect(tag)
                            }
                            }>
                                <Badge variant="secondary" style={{backgroundColor: tag.color}}>{tag.name}</Badge>
                            </CommandItem>
                        )} 
                        </Fragment>
                    ))}

                    </CommandGroup>
                </CommandList>
            </Command>
            <div className="border-t p-2">
                <Button className="w-full" onClick={openCreateMenu}> <Plus className="mr-2 h-4 w-4" /> Create Tag</Button>
            </div>
        </CommandDialog>
        <Dialog open={createOpen} onOpenChange={setCreateOpen}>
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
        </Dialog>
        </>
    )
}