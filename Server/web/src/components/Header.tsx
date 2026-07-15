import {
  Menubar,
//   MenubarContent,
//   MenubarGroup,
//   MenubarItem,
  MenubarMenu,
//   MenubarSeparator,
//   MenubarShortcut,
  MenubarTrigger,
} from "@/components/ui/menubar"

import { House, Dices } from "lucide-react"
import { useNavigate } from "react-router"

export const Header = () => {
    const navigate = useNavigate()

    const goHome = () => {
        navigate("/")
    }

    const goToRandomVideo = () => {
        navigate("video/thisWillBeARandomVideo")
    }

    return (
        <header className="header-container">
            <div className="w-full p-3">
                <Menubar className="h-14">
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
                </Menubar>
            </div>
        </header>
    )
}