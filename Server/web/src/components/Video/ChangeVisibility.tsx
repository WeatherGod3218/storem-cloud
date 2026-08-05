
import { useState } from 'react';

import {
    useMutation,
} from '@tanstack/react-query'

import {
    Command,
    CommandDialog,
    CommandEmpty,
    CommandGroup,
    CommandItem,
    CommandList,
} from "@/components/ui/command"


import {
    Dialog,
    DialogFooter,
    DialogContent,
    DialogDescription,
    DialogHeader,
    DialogTitle,
} from "@/components/ui/dialog"

import { Button } from "@/components/ui/button"

import { useAuth } from '@/context/AuthContext';
import { Eye, EyeOff } from 'lucide-react';

type VisibilityProps = {
    video_id: string
    visibility?: string
    open: boolean
    setOpen: (open: boolean) => void
    onSelect: (visibility: string) => void
}

const ENDPOINT = "/api/v2/videos/visibility";

export const ChangeVisibilityPopup = (props: VisibilityProps) => {
    const { session } = useAuth()
    const [visibility, setVisibility] = useState<string>(props.visibility ? props.visibility : "Private")
    const [publicOpen, setPublicOpen] = useState<boolean>(false)

    const visibilityMutation = useMutation<string, Error, {row_id: string, visibility: string}>({
        mutationKey: [`change-video-visibility`, props.video_id],
        mutationFn: (payload: {row_id: string, visibility: string}) =>
        fetch(`${ENDPOINT}`, {
            method: "PUT",
            headers: {
                "Content-Type": "application/json",
                "Authorization": `Bearer ${session?.access_token}`
            },
            body: JSON.stringify(payload)
        }).then((res) => {
            if (!res.ok) {
                if (res.status == 409) {
                        console.log("unable to do thisss")
                    return
                }
                throw new Error("Error updating visibility")
            }

            return res.json()
        }),
        onSuccess: (newVisibility: string) => {
            setVisibility(newVisibility)
        },
        onSettled: () => {
            exitOut()
        }
    })

    function exitOut() {
        props.setOpen(false)
        setPublicOpen(false)        
    }

    function openPublicMenu() {
        props.setOpen(false)
        setPublicOpen(true)
    }

    function closePublicMenu() {
        props.setOpen(true)
        setPublicOpen(false)
    }

    function setVideoVisibility(newVisibility: string) {
        if (newVisibility === visibility) {
            console.log("doing this!")
            return
        }

        visibilityMutation.mutate({row_id: props.video_id, visibility: newVisibility})
    }

    return (
        <>
        <CommandDialog open={props.open} onOpenChange={props.setOpen}>
            <div className="border-b p-2 text-center">
                <p>Set Video Visibility</p>
            </div>
            <Command>
                <CommandList>
                    <CommandEmpty>No results found.</CommandEmpty>
                    <CommandGroup>
                        <CommandItem value={"Public"} onSelect={() => {
                            openPublicMenu()
                        }
                        }>
                            <div className='flex flex-row items-center'>
                                <Eye color='red'/>
                                <p className='ml-2 text-red-500'>Public</p>
                            </div>
                        </CommandItem>
                        <CommandItem value={"Private"} onSelect={() => {
                            setVideoVisibility("Private")
                            props.onSelect("Private")
                        }
                        }>
                            <div className='flex flex-row items-center'>
                                <EyeOff/>
                                <p className='ml-2'>Private</p>
                            </div>                        
                        </CommandItem>
                    </CommandGroup>
                </CommandList>
            </Command>
        </CommandDialog>
        <Dialog open={publicOpen} onOpenChange={setPublicOpen}>
            <DialogContent className="sm:max-w-sm">        
                    <DialogHeader>
                        <DialogTitle>Are you sure?</DialogTitle>
                        <DialogDescription>
                            Setting this video to public will allow anyone to be able to see it. Make sure there is no questionable content before doing so.
                            If you are unsure, do not make it public.
                        </DialogDescription>
                    </DialogHeader>
                    <DialogFooter className='mt-1'>
                        <Button type="button" variant="outline" onClick={closePublicMenu}>Cancel</Button>
                        <Button type="button" variant="destructive" onClick={() => {
                            setVideoVisibility("Public")
                            props.onSelect("Public")
                        }}>Make Public</Button>
                    </DialogFooter>        
            </DialogContent>
        </Dialog>
        </>
    )
}