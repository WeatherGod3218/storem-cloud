import { Fragment, useCallback, useEffect, useRef, useState } from "react";
import { ScrollBarThumbnail } from "./ScrollBarThumbnail";
import { useAuth } from "@/context/AuthContext";
import { useFilter } from "@/context/FilterContext";

const ENDPOINT = "/api/v2/videos/group";

type Cursor = {
    row_id: string,
    timestamp: number,
}

type Video = {
    row_id: string,
    thumbnail: string,
    filename: string,
    username: string,

	custom_title?: string,
	custom_description?: string,
}

type VideoScrollBarProps = {
	cursor?: Cursor
}

function useVideoGroup(props: VideoScrollBarProps) {
	const {session, authLoading} = useAuth()
	const { filter } = useFilter()

	const [filterElement, setFilterElement] = useState(filter.filter_element)
	const [scrambled, setScrambled] = useState(false)

    const [videos, setVideos] = useState<Video[]>([]);
    const [cursor, setCursor] = useState<Cursor | null>(props.cursor || null);
    const [hasMore, setHasMore] = useState(true);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const hasLoadedOnce = useRef(false);


	function resetVideos(setNull: boolean) {
		setCursor((props.cursor && !setNull) ? props.cursor : null)
		setHasMore(true)
		setVideos([])
	}

  	const loadMore = useCallback(async () => {
		console.log(`Auth: ${authLoading}`)
		console.log(`Has More: ${hasMore}`)
		console.log(`Loading: ${loading}`)
		console.log(`Cursor: ${JSON.stringify(cursor)}`)

		if (loading || !hasMore || authLoading) return;
		setLoading(true);
		setError(null);

		const payload = {"row_id": cursor?.row_id, "timestamp":cursor?.timestamp, "filter": filter}

		try {
			const res = await fetch(`${ENDPOINT}`, {
				method: "POST",
				headers: {
					"Content-Type": `application/json`,
					"Authorization": `Bearer ${session?.access_token}`
				},
				body: JSON.stringify(payload)
			});
			if (!res.ok) throw new Error(`Request failed: ${res.status}`);
			const data: {videos: Video[]; cursor: Cursor | null} = await res.json();

			setVideos((prev) => [...prev, ...(data.videos ?? [])]);
			setCursor(data.cursor ?? null);
			setHasMore(Boolean(data.cursor));
		} catch (err) {
			setError(err instanceof Error ? err.message : "Something went wrong");
			setHasMore(false);
		} finally {
			setLoading(false);
		}
	}, [cursor, hasMore, loading]);

	useEffect(() => {
		if (hasLoadedOnce.current) return;
		hasLoadedOnce.current = true;
		loadMore();
	}, []);

	useEffect(() => {
		if (filter.filter_element != filterElement) {
			resetVideos(true)
		}
		resetVideos(false)
	}, [filter]);

  	return { videos, loadMore, hasMore, loading, error };
}

export const VideoScrollBar = (props: VideoScrollBarProps) => {
    const { videos, loadMore, hasMore, loading, error } = useVideoGroup(props);
    const [isIntersecting, setIsIntersecting] = useState(false);
	const { filter } = useFilter()

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

  return (
    <div className="w-full h-full">
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

		{videos.length === 0 && !loading && !error && (
			<div className="text-center py-16 text-slate-500 text-sm">
				No videos yet.
			</div>
		)}

		<div className="w-full px-3 grid grid-cols-1">
			{videos.map((video) => (
				<Fragment key={video.row_id}>
				{(video.row_id != props.cursor?.row_id) &&
					<ScrollBarThumbnail key={video.row_id} rowId={video.row_id} customTitle={video.custom_title} customDescription={video.custom_description} filename={video.filename} username={video.username} thumbnail={video.thumbnail} />
				}
				</Fragment>
			))}
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
				<span className="text-sm text-slate-400">No More Videos.</span>
			)}
		</div>
	</div>
  );
}
