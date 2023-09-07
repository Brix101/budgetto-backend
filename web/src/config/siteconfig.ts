import { type SidebarNavItem } from "@/types";

export interface DashboardConfig {
  sidebarNav: SidebarNavItem[];
}

export const dashboardConfig: DashboardConfig = {
  sidebarNav: [
    {
      title: "Dashboard",
      href: "/dashboard",
      icon: "terminal",
      items: [],
    },
    {
      title: "Transactions",
      href: "/dashboard/transactions",
      icon: "transaction",
      items: [],
    },
    {
      title: "Budgets",
      href: "/dashboard/budgets",
      icon: "dollarSign",
      items: [],
    },
    {
      title: "Accounts",
      href: "/dashboard/accounts",
      icon: "account",
      items: [],
    },
    {
      title: "Categories",
      href: "/dashboard/categories",
      icon: "category",
      items: [],
    },
  ],
};
