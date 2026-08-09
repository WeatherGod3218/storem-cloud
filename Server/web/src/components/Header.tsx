import {
  Menubar,
  MenubarMenu,
  MenubarTrigger,
} from "@/components/ui/menubar"

import {
	useMutation,
    useQuery,
    useQueryClient,
} from '@tanstack/react-query'

import { useAuth } from "@/context/AuthContext";
import { House, Dices, User, Funnel, ScrollText } from "lucide-react"
import { useNavigate } from "react-router"
import { useEffect, useState } from "react";
import { FilterPopup } from "./Videos/FilterPopup";

const USERS_ENDPOINT = "/api/v2/users/info";
const RANDOOM_VIDEO_ENDPOINT =  "/api/v2/videos/random";

type UserInfo = {
    role: string,
    actions: boolean
}

function toTitleCase(str: string) {
  return str
    .toLowerCase()
    .split(' ')
    .map(word => word.charAt(0).toUpperCase() + word.slice(1))
    .join(' ');
}


export const Header = () => {
    const [filterOpen, setFilterOpen] = useState<boolean>(false)
    const { session, authLoading } = useAuth()
    const navigate = useNavigate()
    const queryClient = useQueryClient()

    const goHome = () => {
        navigate("/")
    }

    const goToLogin = () => {
        navigate("/login")
    }
 
    const goToLogs = () => {
        navigate("/actions")
    }
    
    const { error, data } = useQuery<UserInfo, Error>({
        queryKey: [`get-user-info`],
        retry: 3,
        queryFn: () =>
        fetch(`${USERS_ENDPOINT}`, {
			method: "GET",
			headers: {
				"Authorization": `Bearer ${session?.access_token}`
		    },
        }).then((res) => {
            if (!res.ok) throw new Error(`Failed to get video data: ${res.status}`)
            return res.json()
        }),
    })

    useEffect(() => {
        queryClient.invalidateQueries({
            queryKey: [`get-user-info`]
        })
    }, [authLoading])

    const getRandomVideo = useMutation<any, Error, null>({
		mutationKey: [`get-random-video`],
		mutationFn: () =>
		fetch(`${RANDOOM_VIDEO_ENDPOINT}`, {
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
            <header className="sticky w-full header-container min-w-0">
                <div className="w-full mr-3 min-w-0">
                    <Menubar className="flex h-14 bg-black justify-between min-w-0">
                        <div className="flex flex-row shrink-0 min-w-0">
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
                            { data?.actions &&
                             <MenubarMenu>
                                <MenubarTrigger 
                                className="h-12 text-bold"
                                onClick={goToLogs}
                                ><ScrollText className="h-1/2 mr-1"/>Logs</MenubarTrigger>
                            </MenubarMenu>                                 
                            }  
                        </div>
                        <div className="flex-1 flex justify-end min-w-0">
                            <MenubarMenu>
                                <MenubarTrigger 
                                className="h-12 mr-2 text-bold text-right min-w-0 max-w-full"
                                onClick={goToLogin}
                                >
                                    <User className="h-1/2 shrink-0 mr-1"/> 
                                    {error ? <p className="min-w-0 truncate">Login Failed</p>:
                                        <p className="min-w-0 truncate">{(authLoading ? "Loading" : (data ? toTitleCase(data.role) : "Not Logged In!"))}</p>
                                    }
                                </MenubarTrigger>
                            </MenubarMenu>                         
                        </div>
                    </Menubar>
                </div>
            </header>    
            </nav>    
        </>

    )
}