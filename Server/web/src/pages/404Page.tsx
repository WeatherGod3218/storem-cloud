import { Header } from "@/components/Header"

export const NotFoundPage = () => {

    return (
        <div>
            <title> Test </title>
            <Header/>
            <div className="w-full h-screen flex justify-center">
                <div className="flex flex-col min-h-0 w-full">
                    <div className="mt-5 text-center w-full">
                        <h1 className="pl-1 text-8xl text-bold text-red-100"> 404 </h1>
                        <h2 className="pl-1 text-xl text-slate-200"> Sorry, but we are unable to find the content you are looking for.</h2>
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