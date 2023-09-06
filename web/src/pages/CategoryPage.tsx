import { Button } from "@/components/ui/button";
import { QUERY_CATEGORIES_KEY } from "@/constant/query.constant";
import queryClient from "@/lib/queryClient";
import { Category } from "@/lib/validations/category";
import { useGetCategories } from "@/services/category";
import { useAuth0 } from "@auth0/auth0-react";
import { useQuery } from "@tanstack/react-query";
import { AxiosError } from "axios";
import { Suspense } from "react";
import {
  Await,
  LoaderFunctionArgs,
  useAsyncError,
  useLoaderData,
} from "react-router-dom";

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

function CategoryPage() {
  const { logout } = useAuth0();
  const { promiseData } = useLoaderData() as HomeLoader;

  return (
    <div>
      <Button
        onClick={() =>
          logout({ logoutParams: { returnTo: window.location.origin } })
        }
      >
        Logout
      </Button>
      categories
      <Suspense fallback={<div>Loading...</div>}>
        <Await resolve={promiseData} errorElement={<ErrorBoundary />}>
          <SomeView />
        </Await>
      </Suspense>
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

export default CategoryPage;
