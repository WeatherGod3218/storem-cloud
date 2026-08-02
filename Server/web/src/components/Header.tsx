import {
  Menubar,
  MenubarMenu,
  MenubarTrigger,
} from "@/components/ui/menubar"

import {
	useMutation,
} from '@tanstack/react-query'

import { useAuth } from "@/context/AuthContext";
import { House, Dices, User } from "lucide-react"
import { useNavigate } from "react-router"

const ENDPOINT = "/api/v2/videos/random";

export const Header = () => {
    const { user, session, authLoading } = useAuth()
    const navigate = useNavigate()

    const goHome = () => {
        navigate("/")
    }

    const goToLogin = () => {
        navigate("/login")
    }
    //TE
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
        <header className="header-container">
            <div className="w-full p-3">
                <Menubar className="h-14 justify-between">
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
                    </div>
                    <div>
                        <MenubarMenu>
                            <MenubarTrigger 
                            className="h-12 text-bold"
                            onClick={goToLogin}
                            ><User className="h-1/2 mr-1"/>{(authLoading ? "Loading" : user ? user.email : "Not Logged In!")}</MenubarTrigger>
                        </MenubarMenu>                         
                    </div>
                </Menubar>
            </div>
        </header>
    )
}