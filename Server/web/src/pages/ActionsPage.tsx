import useDocumentTitle from "@/hooks/DocumentTitle"
import { Header } from "../components/Header"
import ActionList from "@/components/Actions/ActionList"
import { useNavigate } from "react-router";
import { useQuery } from "@tanstack/react-query";
import { useAuth } from "@/context/AuthContext";
import { useEffect } from "react";

const ENDPOINT = "/api/v2/users/info";

type UserInfo = {
    role: string,
    actions: boolean
}

class HttpError extends Error {
  status: number;

  constructor(message: string, status: number) {
    super(message);
    this.name = "HttpError";
    this.status = status;
  }
}

export const ActionsPage = () => {
    const { session, authLoading} = useAuth()

    const navigate = useNavigate()

    const goToUnauthorized = () => {
        navigate("/unauthorized")
    }

    useDocumentTitle("Audit Log")

    const { isPending, error } = useQuery<null, Error, UserInfo>({
        queryKey: [`check-action-auth`],
        retry: 3,
        queryFn: () =>
        fetch(`${ENDPOINT}`, {
			method: "GET",
			headers: {
				"Authorization": `Bearer ${session?.access_token}`
		    },
        }).then((res) => {
            if (!res.ok) throw new HttpError(`Failed to get video data: ${res.status}`, res.status)
            return null
        }),
    })

    useEffect(() => {
        if (error instanceof HttpError && error.status === 401) {
            if (session != null) {
                goToUnauthorized();                
            }
        }
    }, [error]);


    return (
        <div>
            <meta name="description" content="Audit Log for Store 'em Cloud"></meta>
            <Header/>
            { !(isPending || authLoading) &&
                <ActionList/>
            }
        </div>
    )
}