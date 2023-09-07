import { Outlet } from "react-router-dom";

import { MainNav } from "@/components/layouts/main-nav";
import { SidebarNav } from "@/components/layouts/sidebar-nav";
import { SiteHeader } from "@/components/layouts/site-header";
import { PageHeader } from "@/components/page-header";
import { Shell } from "@/components/shells/shell";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Skeleton } from "@/components/ui/skeleton";
import { dashboardConfig } from "@/config/siteconfig";
import { useAuth0 } from "@auth0/auth0-react";

export function DashboardShell() {
  const auth = useAuth0();

  if (!auth.isAuthenticated && !auth.isLoading) {
    auth.loginWithRedirect();
    return <DashboardShellLoader />;
  }

  if (auth.isLoading) {
    return <DashboardShellLoader />;
  }

  return (
    <div className="flex min-h-screen flex-col">
      <SiteHeader auth={auth} />
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
  );
}

function DashboardShellLoader() {
  return (
    <div className="flex min-h-screen flex-col">
      <header className="sticky top-0 z-50 w-full border-b bg-background">
        <div className="container flex h-16 items-center">
          <MainNav />
          <div className="flex flex-1 items-center justify-end space-x-4">
            <nav className="flex items-center space-x-2">
              <Skeleton className="h-8 w-8 rounded-full" />
            </nav>
          </div>
        </div>
      </header>
      <div className="container flex-1 items-start md:grid md:grid-cols-[220px_minmax(0,1fr)] md:gap-6 lg:grid-cols-[240px_minmax(0,1fr)] lg:gap-10">
        <aside className="fixed top-14 z-30 -ml-2 hidden h-[calc(100vh-3.5rem)] w-full shrink-0 overflow-y-auto border-r md:sticky md:block">
          <ScrollArea className="py-6 pr-6 lg:py-8">
            <SidebarNav items={dashboardConfig.sidebarNav} isLoading />
          </ScrollArea>
        </aside>
        <main className="flex w-full flex-col overflow-hidden">
          <Shell variant="sidebar">
            <PageHeader
              id="dashboard-stores-page-header"
              aria-labelledby="dashboard-stores-page-header-heading"
            >
              <div className="flex space-x-4 justify-between">
                <Skeleton className="h-10 w-64" />
                <Skeleton className="h-10 w-40" />
              </div>
              <Skeleton className="h-5 w-36" />
            </PageHeader>
            <section></section>
          </Shell>
        </main>
      </div>
    </div>
  );
}
