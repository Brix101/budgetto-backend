import { PageHeader, PageHeaderHeading } from "@/components/page-header";
import { Shell } from "@/components/shells/shell";

function DashboardPage() {
  return (
    <Shell variant="sidebar">
      <PageHeader
        id="dashboard-stores-page-header"
        aria-labelledby="dashboard-stores-page-header-heading"
      >
        <div className="flex space-x-4">
          <PageHeaderHeading size="sm" className="flex-1">
            Dashboard
          </PageHeaderHeading>
        </div>
        {/* <PageHeaderDescription className="" size="sm">
          Manage your dashboard
        </PageHeaderDescription> */}
      </PageHeader>
      <section
        id="dashboard-stores-page-stores"
        aria-labelledby="dashboard-stores-page-stores-heading"
      ></section>
    </Shell>
  );
}

export default DashboardPage;
