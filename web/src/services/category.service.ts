import * as z from "zod";

import { QUERY_CATEGORIES_KEY } from "@/constant/query.constant";
import api from "@/lib/api";
import {
  Category,
  categoriesSchema,
  categorySchema,
  createCategorySchema,
} from "@/lib/validations/category";
import { Auth0ContextInterface, User, useAuth0 } from "@auth0/auth0-react";
import { UseQueryOptions, useQuery } from "@tanstack/react-query";
import { AxiosError } from "axios";

export const getCategories = async (auth: Auth0ContextInterface<User>) => {
  const token = await auth.getAccessTokenSilently();

  const res = await api.get("v1/categories", {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });
  return categoriesSchema.parse({ categories: res.data }).categories;
};

export const useQueryCategories = (
  options?: UseQueryOptions<
    Category[],
    AxiosError,
    Category[],
    readonly [string]
  >
) => {
  const auth = useAuth0();

  return useQuery({
    queryKey: [QUERY_CATEGORIES_KEY],
    queryFn: () => getCategories(auth),
    ...options,
  });
};

export const createCategory = async ({
  auth,
  category,
}: {
  auth: Auth0ContextInterface<User>;
  category: z.infer<typeof createCategorySchema>;
}) => {
  const token = await auth.getAccessTokenSilently();

  const res = await api.post("v1/categories", JSON.stringify(category), {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });

  return categorySchema.parse(res.data);
};
