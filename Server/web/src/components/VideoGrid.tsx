import { useCallback, useEffect, useRef, useState } from "react";
import {
	useMutation,
} from '@tanstack/react-query'

import { Button } from "@/components/ui/button"
import { ArrowUpAZ, ArrowUpZA } from "lucide-react"

import { ThumbnailCard, ThumbnailSkeletonCard } from "./Thumbnails/ThumbnailCard";
import { useAuth } from "@/context/AuthContext";


const ENDPOINT = "/api/v2/videos/group";

type Cursor = {
    row_id: string,
    timestamp: string,
}

type Video = {
    row_id: string,
    thumbnail: string,
    filename: string,
    username: string,

	custom_title?: string,
	custom_description?: string,
	visibility?: string,

	timestamp: string,
}

type VideoData = {
	videos: Video[]; 
	cursor: Cursor | null
}

type VideoFetch = {
	row_id?: string, 
	timestamp?: string,
	order_ascending: boolean
}

function useVideoGroup() {
	const { session, authLoading} = useAuth()
    const [videos, setVideos] = useState<Video[]>([]);
    const [cursor, setCursor] = useState<Cursor | null>(null);
    const [hasMore, setHasMore] = useState(true);
    const [loading, setLoading] = useState(false);
	const [newestFirst, setNewestFirst] = useState(true)
    const [error, setError] = useState<string | null>(null);
    const hasLoadedOnce = useRef(false);


	const getVideos = useMutation<VideoData, Error, VideoFetch>({
		mutationKey: [`get-video-grid-list`],
		mutationFn: (payload: any) =>
		fetch(`${ENDPOINT}`, {
			method: "POST",
			body: JSON.stringify(payload),
			headers: {
				"Content-Type": "application/json",
				"Authorization": `Bearer ${session?.access_token}`
			},
		}).then((res) => {
			if (!res.ok) throw new Error(`Failed to fetch tags: ${res.status}`)
			return res.json()
		}),

		onError: (err) => {
			setError(err instanceof Error ? err.message : "Something went wrong");
			setHasMore(false);			
		},
		onSuccess: (data) => {
			setVideos((prev) => [...prev, ...(data.videos ?? [])]);
			setCursor(data.cursor ?? null);
			setHasMore(Boolean(data.cursor));			
		},
		onSettled: () => {
			setLoading(false);
		}
	})

	function resetVideos() {
		setCursor(null)
		setHasMore(true)
		setVideos([])
	}

  	const loadMore = useCallback(async () => {
		if (loading || authLoading || !hasMore) return;
		setLoading(true);
		setError(null);

		const payload = {
			"row_id": cursor?.row_id, 
			"timestamp": cursor?.timestamp,
			"order_ascending": !newestFirst, //oops
		}

		getVideos.mutate(payload)
	}, [cursor, hasMore, loading, authLoading]);

	useEffect(() => {
		if (hasLoadedOnce.current) return;
		hasLoadedOnce.current = true;
		loadMore();
	}, []);

  	return { videos, newestFirst, setNewestFirst, loadMore, hasMore, loading, error, resetVideos};
}

export default function VideoGridInfinite() {
    const { authLoading} = useAuth()
	const { videos, loadMore, hasMore, loading, newestFirst, setNewestFirst, error, resetVideos } = useVideoGroup();
    const [isIntersecting, setIsIntersecting] = useState(false);

    const sentinelRef = useRef(null);

    const loadMoreRef = useRef(loadMore);
    useEffect(() => {
    	loadMoreRef.current = loadMore;
    }, [loadMore]);

    useEffect(() => {
        if (isIntersecting && hasMore && !loading) {
            loadMore();
        }
    }, [isIntersecting, hasMore, loading, loadMore]);

   useEffect(() => {
        const el = sentinelRef.current;
        if (!el) return;

        const observer = new IntersectionObserver(
            (entries) => setIsIntersecting(entries[0].isIntersecting),
            { rootMargin: "400px" }
        );

        observer.observe(el);
        return () => observer.disconnect();
    }, []);


	function changeAscensionOrder() {
		setNewestFirst(!newestFirst)		
		resetVideos()
	}

  	return (
    <div className="min-h-screen bg-dark p-8">
    	<div className="max-w-5xl mx-auto">
			<div className="mb-6">
				<h1 className="text-xl font-semibold text-slate-200">Videos</h1>
				<div className="flex flex-row items-center">					
					<Button variant="outline" size="icon-sm" aria-label="Confirm Title" onClick={changeAscensionOrder}>{ newestFirst ? <ArrowUpAZ/> : <ArrowUpZA/>}</Button>
					<p className="pl-1 text-l text-slate-500">{ newestFirst ? "Newest First" : "Oldest First"}</p>
				</div>
			</div>
			{error && (
				<div className="mb-6 flex items-center gap-2 rounded-lg border border-red-200 bg-orange-950 px-4 py-3 text-sm text-red-700">
					<span>Couldn't load videos: {error}</span>
					<button
					onClick={loadMore}
					className="ml-auto text-red-700 underline underline-offset-2 hover:text-red-800"
					>
					Retry
					</button>
				</div>
			)}
			{error && (
				<div className="mb-6 flex items-center gap-2 rounded-lg border border-red-200 bg-orange-950 px-4 py-3 text-sm text-red-700">
					<span>Couldn't load videos: {error}</span>
					<button
					onClick={loadMore}
					className="ml-auto text-red-700 underline underline-offset-2 hover:text-red-800"
					>
					Retry
					</button>
				</div>
			)}

			{ (authLoading) &&
				<div className="text-center py-16 text-slate-500 text-sm">
					Loading Videos. Please Wait.
				</div>				
			}
			{videos.length === 0 && !authLoading && !loading && !error && (
				<div className="text-center py-16 text-slate-500 text-sm">
					No videos yet.
				</div>
			)}

			<div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
				{videos.map((video) => (
					<ThumbnailCard key={video.row_id} {...video} />
				))}
				{loading && Array.from({ length: 6 }).map((_, i) => <ThumbnailSkeletonCard key={`sk-${i}`} />)}
			</div>

			<div ref={sentinelRef} className="flex justify-center py-8">
				{!loading && hasMore && (
				<button
					onClick={loadMore}
					className="inline-flex items-center gap-2 rounded-lg border border-slate-200 bg-white px-4 py-2 text-sm font-medium text-slate-700 hover:bg-slate-100 transition-colors"
				>
				Load more
				</button>
				)}
				{loading && (
					<span className="inline-flex items-center gap-2 text-sm text-slate-500">
					Loading more videos…
					</span>
				)}
				{!hasMore && videos.length > 0 && (
					<span className="text-sm text-slate-400">You've reached the end.</span>
				)}
			</div>
    	</div>
	</div>
  );
}
