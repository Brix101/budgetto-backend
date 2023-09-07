
import { DataTable } from "@/components/data-table";
import { categoryColumns } from "@/components/data-table/colums";
import { CategoryCreateDialog } from "@/components/forms/category-create";
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
        id="dashboard-stores-page-header"
        aria-labelledby="dashboard-stores-page-header-heading"
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
        id="dashboard-stores-page-stores"
        aria-labelledby="dashboard-stores-page-stores-heading"
      >
        <DataTable data={data ?? []} columns={categoryColumns} />
      </section>
    </Shell>
  );
}


export default CategoryPage;
