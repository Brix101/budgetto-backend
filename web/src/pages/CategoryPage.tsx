
import { DataTable } from "@/components/data-table";
import { categoryColumns } from "@/components/data-table/colums";
import { CategoryCreateDialog } from "@/components/forms/category-form";
import {
  PageHeader,
  PageHeaderDescription,
  PageHeaderHeading,
} from "@/components/page-header";
import { Shell } from "@/components/shells/shell";
import { useQueryCategories } from "@/services/category.service";


function CategoryPage() {
  const { data } = useQueryCategories();

  return (
    <Shell variant="sidebar">
      <PageHeader
        id="dashboard-categories-page-header"
        aria-labelledby="dashboard-categories-page-header-heading"
      >
        <div className="flex space-x-4">
          <PageHeaderHeading size="sm" className="flex-1">
            Categories
          </PageHeaderHeading>
          <CategoryCreateDialog />
        </div>
        <PageHeaderDescription size="sm">
          Manage your categories
        </PageHeaderDescription>
      </PageHeader>
      <section
        id="dashboard-categories-page-categories"
        aria-labelledby="dashboard-categories-page-categories-heading"
      >
        <DataTable data={data ?? []} columns={categoryColumns} searchPlaceHolder="Filter categories..."/>
      </section>
    </Shell>
  );
}


export default CategoryPage;
