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

const ENDPOINT = "/api/v2/videos/random";

export const Header = () => {
    const navigate = useNavigate()

    const goHome = () => {
        navigate("/")
    }

    const goToRandomVideo = () => {
        let cancelled = false;

        fetch(`${ENDPOINT}`)
        .then(res => {if (!res.ok) {throw new Error(`HTTP ERROR: ${res.status}`)} return res.json()})
        .then(data => { if (!cancelled)  navigate(`/video/${data.row_id}`); })        
        .catch(err => { if (!cancelled) {console.log(`${err}`)}})
        .finally(() => { if (!cancelled) cancelled = true; })
        return () => { cancelled = true; };
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