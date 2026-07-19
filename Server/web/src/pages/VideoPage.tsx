import { Header } from "../components/Header"

import { useState, useEffect } from "react"
import { useParams } from "react-router"
import { VideoCard, SkeletonVideoCard } from "@/components/Video/VideoCard"
//import { useMediaQuery } from "@/hooks/MediaQuery";

const ENDPOINT = "/api/v2/videos/video";

type VideoData = {
    row_id: string,
    s3_id: string,

    custom_title: string | null,
    custom_description: string | null,

    username: string,
    filename: string,
    video_url: string,
}

export const VideoPage = () => {    
    let params = useParams()
    //const isVertical = useMediaQuery("(max-width: 768px)");

    const [isLoading, setLoading] = useState<boolean>(true)
    const [data, setData] = useState<VideoData>({
        row_id: "",
        s3_id: "",
        custom_title: "",
        custom_description: "",
        username: "",
        filename: "",
        video_url: "",
    })


    useEffect(() => {
        let cancelled = false
        setLoading(true);

        fetch(`${ENDPOINT}/${params.id}`)
        .then(res => {if (!res.ok) {throw new Error(`HTTP ERROR: ${res.status}`)} return res.json()})
        .then(data => { if (!cancelled) setData(data); })        
        .catch(err => { if (!cancelled) {console.log(`${err}`)}})
        .finally(() => { if (!cancelled) setLoading(false); })
        return () => { cancelled = true; };
    }, [params.id]);

    console.log(params.id)

    return (
        <div className="w-full h-screen flex flex-col">
            <Header/>
            <div className="flex-1 w-full min-h-0 overflow-hidden">
                {!isLoading ? <VideoCard {...data}
                /> : <SkeletonVideoCard/>}
            </div>
        </div>
    )
}