import { DataTable } from "@/components/data-table";
import { categoryColumns } from "@/components/data-table/colums";
import {
  CategoryCreateDialog,
  CategoryDeleteDialog,
  CategoryUpdateDialog,
} from "@/components/forms/category-form";
import {
  PageHeader,
  PageHeaderDescription,
  PageHeaderHeading,
} from "@/components/page-header";
import { Shell } from "@/components/shells/shell";
import { Button } from "@/components/ui/button";
import { useBoundStore } from "@/lib/store";
import { useQueryCategories } from "@/services/category.service";

function CategoryPage() {
  const { setMode } = useBoundStore((state) => state.category);
  const { data } = useQueryCategories();

  function handleCreateClick() {
    setMode({ mode: "create" });
  }

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
          <Button size="sm" onClick={handleCreateClick}>
            Create category
          </Button>
        </div>
        <PageHeaderDescription size="sm">
          Manage your categories
        </PageHeaderDescription>
      </PageHeader>
      <section
        id="dashboard-categories-page-categories"
        aria-labelledby="dashboard-categories-page-categories-heading"
      >
        <DataTable
          data={data ?? []}
          columns={categoryColumns}
          searchPlaceHolder="Filter categories..."
        />

        <CategoryCreateDialog />
        <CategoryDeleteDialog />
        <CategoryUpdateDialog />
      </section>
    </Shell>
  );
}

export default CategoryPage;
