import useDocumentTitle from "@/hooks/DocumentTitle"
import { Header } from "../components/Header"
import VideoGridInfinite from "@/components/Videos/VideoGrid"

export const LandingPage = () => {
    useDocumentTitle("Store 'em Cloud")

    return (
        <div>
            <meta name="description" content="Authorized Video Player for an Automated CLOUD (get it?) Video Backup"></meta>
            <Header/>
            <VideoGridInfinite/>
        </div>
    )
}