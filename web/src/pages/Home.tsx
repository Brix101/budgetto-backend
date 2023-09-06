import { Button } from "@/components/ui/button";
import { QUERY_CATEGORIES_KEY } from "@/constant/query.constant";
import queryClient from "@/lib/queryClient";
import { Category } from "@/lib/validations/category";
import { useGetCategories } from "@/services/category";
import { useAuth0 } from "@auth0/auth0-react";
import { useQuery } from "@tanstack/react-query";
import { AxiosError } from "axios";
import { LoaderFunctionArgs, useAsyncError } from "react-router-dom";

type HomeLoader = {
  promiseData: Promise<Category[]>;
};

export const loader = ({ request, params }: LoaderFunctionArgs): HomeLoader => {
  console.log({ request, params });
  return {
    promiseData: queryClient.fetchQuery(
      [QUERY_CATEGORIES_KEY],
      useGetCategories,
      {
        staleTime: 10000,
      }
    ),
  };
};

function Home() {
  const {
    loginWithPopup,
    isAuthenticated,
    logout,
    user,
    isLoading,
    getAccessTokenSilently,
    getIdTokenClaims,
  } = useAuth0();

  if (isLoading) {
    return <div>Loading ...</div>;
  }

  getIdTokenClaims().then((res) => console.log(res));
  getAccessTokenSilently().then((res) => console.log({ res }));
  return (
    <div>
      {isAuthenticated && user ? (
        <div>
          <img src={user.picture} alt={user.name} />
          <h2>{user.name}</h2>
          <p>{user.email}</p>
        </div>
      ) : undefined}
      {isAuthenticated ? (
        <Button
          onClick={() =>
            logout({ logoutParams: { returnTo: window.location.origin } })
          }
        >
          Logout
        </Button>
      ) : (
        <Button onClick={() => loginWithPopup()}>Login</Button>
      )}
    </div>
  );
}

function ErrorBoundary() {
  const error = useAsyncError() as AxiosError;
  console.log({ error });
  return <div>Dang! {error.message}</div>;
}

function SomeView() {
  const { data } = useQuery<Category[]>([QUERY_CATEGORIES_KEY]);

  return (
    <div>
      {data?.map((category) => {
        return <div key={category.id}>{category.name}</div>;
      })}
    </div>
  );
}

export default Home;
