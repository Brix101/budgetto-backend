import { type SidebarNavItem } from "@/types";
import { useAuth0 } from "@auth0/auth0-react";
import { Outlet } from "react-router-dom";
import { ScrollArea } from "../ui/scroll-area";
import { SidebarNav } from "./sidebar-nav";
import { SiteHeader } from "./site-header";

export interface DashboardConfig {
  sidebarNav: SidebarNavItem[];
}

export const dashboardConfig: DashboardConfig = {
  sidebarNav: [
    {
      title: "Categories",
      href: "/dashboard/categories",
      icon: "user",
      items: [],
    },
    {
      title: "Accounts",
      href: "/dashboard/accounts",
      icon: "store",
      items: [],
    },
    {
      title: "Budgets",
      href: "/dashboard/budgets",
      icon: "dollarSign",
      items: [],
    },
    {
      title: "Transactions",
      href: "/dashboard/transactions",
      icon: "billing",
      items: [],
    },
  ],
};

function DashboardLayout() {
  const auth = useAuth0();

  if (!auth.isAuthenticated && !auth.isLoading) {
    auth.loginWithRedirect();
    return <div>Loading ...</div>;
  }

  if (auth.isLoading) {
    return <div>Loading ...</div>;
  }

  return (
    <>
      <div className="flex min-h-screen flex-col">
        <SiteHeader />
        <div className="container flex-1 items-start md:grid md:grid-cols-[220px_minmax(0,1fr)] md:gap-6 lg:grid-cols-[240px_minmax(0,1fr)] lg:gap-10">
          <aside className="fixed top-14 z-30 -ml-2 hidden h-[calc(100vh-3.5rem)] w-full shrink-0 overflow-y-auto border-r md:sticky md:block">
            <ScrollArea className="py-6 pr-6 lg:py-8">
              <SidebarNav items={dashboardConfig.sidebarNav} className="p-1" />
            </ScrollArea>
          </aside>
          <main className="flex w-full flex-col overflow-hidden">
            <Outlet />
          </main>
        </div>
      </div>
    </>
  );
}

export default DashboardLayout;
