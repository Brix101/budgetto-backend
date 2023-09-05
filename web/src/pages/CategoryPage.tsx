import { QUERY_CATEGORIES_KEY } from "@/constant/query.constant";
import { Category } from "@/lib/validations/category";
import { useQuery } from "@tanstack/react-query";

function CategoryPage() {
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
