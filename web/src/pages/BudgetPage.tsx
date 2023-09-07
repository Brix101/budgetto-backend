import {
  PageHeader,
  PageHeaderDescription,
  PageHeaderHeading,
} from "@/components/page-header";
import { Shell } from "@/components/shells/shell";
import { buttonVariants } from "@/components/ui/button";
import { cn } from "@/lib/utils";
import { Link } from "react-router-dom";

function BudgetPage() {
  return (
    <Shell variant="sidebar">
      <PageHeader
        id="dashboard-stores-page-header"
        aria-labelledby="dashboard-stores-page-header-heading"
      >
        <div className="flex space-x-4">
          <PageHeaderHeading size="sm" className="flex-1">
            Budgets
          </PageHeaderHeading>
          <Link
            aria-label="Create budget"
            to={"/dashboard/budgets"}
            // href={getDashboardRedirectPath({
            //   storeCount: allStores.length,
            //   subscriptionPlan: subscriptionPlan,
            // })}
            className={cn(
              buttonVariants({
                size: "sm",
              })
            )}
          >
            Create budget
          </Link>
        </div>
        <PageHeaderDescription size="sm">
          Manage your budgets
        </PageHeaderDescription>
      </PageHeader>
      <section
        id="dashboard-stores-page-stores"
        aria-labelledby="dashboard-stores-page-stores-heading"
      ></section>
    </Shell>
  );
}

export default BudgetPage;
