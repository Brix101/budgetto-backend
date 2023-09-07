import {
  PageHeader,
  PageHeaderDescription,
  PageHeaderHeading,
} from "@/components/page-header";
import { Shell } from "@/components/shells/shell";
import { buttonVariants } from "@/components/ui/button";
import { cn } from "@/lib/utils";
import { Link } from "react-router-dom";

function TransactionPage() {
  return (
    <Shell variant="sidebar">
      <PageHeader
        id="dashboard-stores-page-header"
        aria-labelledby="dashboard-stores-page-header-heading"
      >
        <div className="flex space-x-4">
          <PageHeaderHeading size="sm" className="flex-1">
            Transactions
          </PageHeaderHeading>
          <Link
            aria-label="Create transaction"
            to={"/dashboard/transactions"}
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
            Create transaction
          </Link>
        </div>
        <PageHeaderDescription size="sm">
          Manage your transactions
        </PageHeaderDescription>
      </PageHeader>
      <section
        id="dashboard-stores-page-stores"
        aria-labelledby="dashboard-stores-page-stores-heading"
      ></section>
    </Shell>
  );
}

export default TransactionPage;
