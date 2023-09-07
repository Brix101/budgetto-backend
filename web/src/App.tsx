import { Toaster } from "@/components/ui/toaster";
import AccountPage from "@/pages/AccountPage";
import BudgetPage from "@/pages/BudgetPage";
import CategoryPage from "@/pages/CategoryPage";
import DashboardPage from "@/pages/DashboardPage";
import Home from "@/pages/Home";
import PageNotFound from "@/pages/PageNotFound";
import TransactionPage from "@/pages/TransactionPage";
import { ReactQueryDevtools } from "@tanstack/react-query-devtools";
import { createBrowserRouter, RouterProvider } from "react-router-dom";
import { DashboardShell } from "./components/shells/layout-shell";

function App() {
  const router = createBrowserRouter([
    {
      path: "/",
      element: <Home />,
    },
    {
      path: "/dashboard",
      element: <DashboardShell />,
      children: [
        {
          index: true,
          element: <DashboardPage />,
        },
        {
          path: "accounts",
          element: <AccountPage />,
        },
        {
          path: "budgets",
          element: <BudgetPage />,
        },
        {
          path: "categories",
          element: <CategoryPage />,
        },
        {
          path: "transactions",
          element: <TransactionPage />,
        },
      ],
    },
    {
      path: "*",
      element: <PageNotFound />,
    },
  ]);

  return (
    <>
      <RouterProvider router={router} fallbackElement={<div>Loading...</div>} />
      <ReactQueryDevtools initialIsOpen={false} position="bottom-right" />
      <Toaster />
    </>
  );
}

export default App;
