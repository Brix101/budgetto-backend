import { QUERY_CATEGORIES_KEY } from "@/constant/query.constant";
import api from "@/lib/api";
import { Category, categoriesSchema } from "@/lib/validations/category";
import { UseQueryOptions, useQuery } from "@tanstack/react-query";
import { AxiosError } from "axios";

export const useGetCategories = async () => {
  const res = await api.get("v1/categories");
  return categoriesSchema.parse({ categories: res.data }).categories;
};

export const useQueryCategories = (
  options?: UseQueryOptions<
    Category[],
    AxiosError,
    Category[],
    readonly [string]
  >,
) => {
  return useQuery({
    queryKey: [QUERY_CATEGORIES_KEY],
    queryFn: useGetCategories,
    ...options,
  });
};
