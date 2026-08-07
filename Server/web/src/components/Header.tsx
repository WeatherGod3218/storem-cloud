import {
  Menubar,
  MenubarMenu,
  MenubarTrigger,
} from "@/components/ui/menubar"

import {
	useMutation,
} from '@tanstack/react-query'

import { useAuth } from "@/context/AuthContext";
import { House, Dices, User, Funnel } from "lucide-react"
import { useNavigate } from "react-router"
import { useState } from "react";
import { FilterPopup } from "./FilterPopup";

const ENDPOINT = "/api/v2/videos/random";

export const Header = () => {
    const [filterOpen, setFilterOpen] = useState<boolean>(false)
    const { user, session, authLoading } = useAuth()
    const navigate = useNavigate()

    const goHome = () => {
        navigate("/")
    }

    const goToLogin = () => {
        navigate("/login")
    }
    
    const getRandomVideo = useMutation<any, Error, null>({
		mutationKey: [`get-random-video`],
		mutationFn: () =>
		fetch(`${ENDPOINT}`, {
			method: "GET",
			headers: {
				"Content-Type": "application/json",
				"Authorization": `Bearer ${session?.access_token}`
			},
		}).then((res) => {
			if (!res.ok) throw new Error(`Failed to fetch tags: ${res.status}`)
			return res.json()
		}),

		onSuccess: (data) => {
			navigate(`/video/${data.row_id}`)			
		},

	})

    const goToRandomVideo = () => {
        getRandomVideo.mutate(null)
    }


    return (
        <>
            <FilterPopup tags={[]} open={filterOpen} setOpen={setFilterOpen}/>
            <nav className="sticky w-full top-0 z-50 flex shadow-sm">
            <header className="sticky w-full header-container">
                <div className="w-full mr-3">
                    <Menubar className="flex h-14 bg-black justify-between">
                        <div className="flex flex-row">
                            <MenubarMenu>
                                <MenubarTrigger 
                                className="h-12 text-bold"
                                onClick={goHome}
                                ><House className="h-1/2 mr-1"/>Home</MenubarTrigger>
                            </MenubarMenu>
                            <MenubarMenu>
                                <MenubarTrigger 
                                className="h-12 text-bold"
                                onClick={goToRandomVideo}
                                ><Dices className="h-1/2 mr-1"/>I'm Feeling Lucky</MenubarTrigger>
                            </MenubarMenu>
                            <MenubarMenu>
                                <MenubarTrigger 
                                className="h-12 text-bold"
                                onClick={() => {setFilterOpen(true)}}
                                ><Funnel className="h-1/2 mr-1"/>Filter</MenubarTrigger>
                            </MenubarMenu>    
                        </div>
                        <div className="flex-1 flex justify-end">
                            <MenubarMenu>
                                <MenubarTrigger 
                                className="h-12 text-bold text-right"
                                onClick={goToLogin}
                                ><User className="h-1/2 mr-1 shrink-0 truncate"/> <p className="min-w-0 truncate">{(authLoading ? "Loading" : user ? user.email : "Not Logged In!")}</p></MenubarTrigger>
                            </MenubarMenu>                         
                        </div>
                    </Menubar>
                </div>
            </header>    
            </nav>    
        </>

    )
}