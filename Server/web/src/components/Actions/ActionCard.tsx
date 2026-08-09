import {
    Card,
    CardContent,   
    CardFooter,
    CardHeader,
    CardTitle,
} from "@/components/ui/card"

import { cn } from "@/lib/utils"
import { Skeleton } from "../ui/skeleton"

type ThumbnailCardProps = {
    row_id: string,
    user_id: string,
    user_email: string
    action: string,
    timestamp: number
}

export const ActionCard = (props: ThumbnailCardProps) => {
    return (
    <Card 
        role="banner"
        tabIndex={0}
        className={cn(
                "relative mx-auto w-full pt-0"
            )}
        >
        <CardHeader>
            <CardTitle className="mt-2 w-full text-center min-w-0 text-2xl"><h1>{props.action}</h1></CardTitle>
        </CardHeader>
        <CardContent>
            <div className="flex flex-row justify-center">
                <> 
                <h2>{props.user_email} (<i>{props.user_id}</i>)</h2>
                </>
            </div>
        </CardContent>
        <div className="flex-1"/> 
        <CardFooter>
            <h4> <i>Time: {new Date(props.timestamp * 1000).toString()} </i></h4>
        </CardFooter>
    </Card>
    )
}

export const SkeletonActionCard = () => {
    return (
    <Card 
        role="banner"
        tabIndex={0}
        className={cn(
                "relative mx-auto w-full pt-0"
            )}
        >
        <CardHeader className="flex justify-center items-center">
            <Skeleton className="mt-3 w-3/4 h-[40px]"/>
        </CardHeader>
        <CardContent>
            <div className="flex flex-row justify-center">
                <> 
                <Skeleton className="w-1/3 h-[30px]"/>
                </>
            </div>
        </CardContent>
        <div className="flex-1"/> 
        <CardFooter>
            <Skeleton className="w-1/4 h-[25px]"/>
        </CardFooter>
    </Card>
    )
}
