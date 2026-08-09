import { useCallback, useEffect, useRef, useState } from "react";
import {
	useMutation,
} from '@tanstack/react-query'

import { useAuth } from "@/context/AuthContext";
import { ActionCard, SkeletonActionCard } from "./ActionCard";


const ENDPOINT = "/api/v2/actions/group";

type Cursor = {
    row_id: string,
    timestamp: number,
}

type Action = {
    row_id: string,
    user_id: string,
	user_email: string,
    
	action: string,
	timestamp: number,
}

type ActionData = {
	actions: Action[]; 
	cursor: Cursor | null
}

type ActionFetch = {
	row_id?: string, 
	timestamp?: number,
}

function useActionGroup() {
	const { session, authLoading} = useAuth()
    const [actions, setActions] = useState<Action[]>([]);
    const [cursor, setCursor] = useState<Cursor | null>(null);
    const [hasMore, setHasMore] = useState(true);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const hasLoadedOnce = useRef(false);


	const getActions = useMutation<ActionData, Error, ActionFetch>({
		mutationKey: [`get-actions-list`],
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
			setActions((prev) => [...prev, ...(data.actions ?? [])]);
			setCursor(data.cursor ?? null);
			setHasMore(Boolean(data.cursor));			
		},
		onSettled: () => {
			setLoading(false);
		}
	})

  	const loadMore = useCallback(async () => {
		if (loading || authLoading || !hasMore) return;
		setLoading(true);
		setError(null);

		const payload = {
			"row_id": cursor?.row_id, 
			"timestamp": cursor?.timestamp,
		}

		getActions.mutate(payload)
	}, [cursor, hasMore, loading, authLoading]);


	useEffect(() => {
		if (hasLoadedOnce.current) return;
		hasLoadedOnce.current = true;
		loadMore();
	}, []);

  	return { actions, loadMore, hasMore, loading, error};
}

export default function ActionList() {
    const { authLoading } = useAuth()

	const { actions, loadMore, hasMore, loading, error } = useActionGroup();
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
    	<div className="w-full">
			<div className="mb-6">
				<h1 className="text-xl font-semibold text-slate-200">Recent Actions</h1>
				<div className="flex flex-row items-center">					
					<p className="pl-1 text-l text-slate-500"> Newest First </p>
				</div>
			</div>
			{error && (
				<div className="mb-6 flex items-center gap-2 rounded-lg border border-red-200 bg-orange-950 px-4 py-3 text-sm text-red-700">
					<span>Couldn't load audit log: {error}</span>
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
					Loading Audit Log. Please Wait.
				</div>				
			}
			{actions.length === 0 && !authLoading && !loading && !error && (
				<div className="text-center py-16 text-slate-500 text-sm">
					No Actions to Load.
				</div>
			)}

			<div className="flex flex-col w-full gap-4">
				{actions.map((action) => (
					<ActionCard key={action.row_id} {...action} />
				))}
				{loading && Array.from({ length: 6 }).map((_, i) => <SkeletonActionCard key={`sk-${i}`} />)}
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
					Loading more actions
					</span>
				)}
				{!hasMore && actions.length > 0 && (
					<span className="text-sm text-slate-400">You have reached the end of the Audit Log.</span>
				)}
			</div>
    	</div>
	</div>
  );
}
