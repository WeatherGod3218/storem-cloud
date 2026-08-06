import { Header } from "@/components/Header"
import useDocumentTitle from "@/hooks/DocumentTitle"

export const UnauthorizedPage = () => {
    useDocumentTitle("UNAUTHORIZED")

    return (            
        <div>
            <Header/>
            <div className="w-full h-screen flex justify-center">
                <div className="flex flex-col min-h-0 w-full">
                    <div className="mt-5 text-center w-full">
                        <h1 className="pl-1 text-5xl text-bold text-red-500"> UNAUTHORIZED </h1>
                        <h2 className="pl-1 text-xl text-slate-200"> Sorry, but the content you are trying to access is private.</h2>
                        <p className="pl-1 text-l text-slate-500"> <i>Head back home using the navbar.</i></p>
                    </div>
                    <div className="flex flex-1 items-center justify-center text-center">
                        {/* <p> meow</p> */}
                    </div>                    
                </div>
            </div>
        </div>
    )
}