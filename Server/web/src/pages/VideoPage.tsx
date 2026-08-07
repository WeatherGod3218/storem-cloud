import { Header } from "../components/Header"

import { useQuery } from "@tanstack/react-query"

import { useParams } from "react-router"
import { VideoCard, SkeletonVideoCard } from "@/components/Video/VideoCard"
import { useAuth } from "@/context/AuthContext"
//import { useMediaQuery } from "@/hooks/MediaQuery";
import { useNavigate } from "react-router"
import { useEffect } from "react"
import useDocumentTitle from "@/hooks/DocumentTitle"

const ENDPOINT = "/api/v2/videos/video";

type VideoData = {
    row_id: string,
    s3_id: string,

    custom_title?: string | null,
    custom_description?: string | null,

    visibility: string,

    username: string,
    filename: string,

    video_url: string,
    thumbnail_url: string,
    
    timestamp: number,
    can_modify: boolean
}

//ig bro
class HttpError extends Error {
  status: number;

  constructor(message: string, status: number) {
    super(message);
    this.name = "HttpError";
    this.status = status;
  }
}

export const VideoPage = () => {   
	const { session, authLoading} = useAuth()
    let params = useParams()

    const navigate = useNavigate()

    const goToUnauthorized = () => {
        navigate("/unauthorized")
    }


    const { isPending, error, data } = useQuery<VideoData, Error>({
        queryKey: [`get-video-data`, params.id],
        retry: 3,
        queryFn: () =>
        fetch(`${ENDPOINT}/${params.id}`, {
			method: "GET",
			headers: {
				"Authorization": `Bearer ${session?.access_token}`
		    },
        }).then((res) => {
            if (!res.ok) throw new HttpError(`Failed to get video data: ${res.status}`, res.status)
            return res.json()
        }),
    })

    useEffect(() => {
        if (error instanceof HttpError && error.status === 401) {
            if (session != null) {
                goToUnauthorized();                
            }
        }
    }, [error]);

    useDocumentTitle((data ? (data.custom_title ? data.custom_title : data.filename) : "Loading..."))

    return (
        <div className="w-full h-screen flex flex-col">
            <Header/>
            <div className="flex-1 w-full min-h-0 overflow-hidden">
                {isPending || authLoading ? (
                    <SkeletonVideoCard />
                ) : error ? (
                    <div>Failed to load video.</div>
                ) : !data ? (
                    <div>No video found.</div>
                ) : (
                    <VideoCard {...data} />
                )}
            </div>
        </div>
    )
}