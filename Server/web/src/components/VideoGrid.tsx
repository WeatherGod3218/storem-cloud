import { useCallback, useEffect, useRef, useState } from "react";
import { ThumbnailCard, ThumbnailSkeletonCard } from "./ThumbnailCard";

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
}

function useVideoGroup() {
    const [videos, setVideos] = useState<Video[]>([]);
    const [cursor, setCursor] = useState<Cursor | null>(null); // null = first page
    const [hasMore, setHasMore] = useState(true);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const hasLoadedOnce = useRef(false);

  	const loadMore = useCallback(async () => {
		if (loading || !hasMore) return;
		setLoading(true);
		setError(null);

		const payload = {"row_id": cursor?.row_id, "timestamp":cursor?.timestamp}
		try {
			const res = await fetch(`${ENDPOINT}`, {
				method: "POST",
				headers:{
					"Content-Type": `application/json`
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
		// eslint-disable-next-line react-hooks/exhaustive-deps
	}, []);

  	return { videos, loadMore, hasMore, loading, error };
}

export default function VideoGridInfinite() {
    const { videos, loadMore, hasMore, loading, error } = useVideoGroup();
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

  return (
    <div className="min-h-screen bg-dark p-8">
    	<div className="max-w-5xl mx-auto">
			<div className="mb-6">
				<h1 className="text-xl font-semibold text-slate-200">Videos</h1>
				<p className="text-sm text-slate-500">Newest first</p>
			</div>

			{error && (
				<div className="mb-6 flex items-center gap-2 rounded-lg border border-red-200 bg-red-850 px-4 py-3 text-sm text-red-700">
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

			<div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
				{videos.map((video) => (
					<ThumbnailCard key={video.row_id} rowId={video.row_id} filename={video.filename} username={video.username} thumbnail={video.thumbnail} />
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
