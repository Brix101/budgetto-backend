import DashboardLayout from "@/components/layouts/DashboardLayout";
import CategoryPage, { loader as categoryLoader } from "@/pages/CategoryPage";
import Dashboard from "@/pages/Dashboard";
import Home from "@/pages/Home";
import PageError from "@/pages/PageError";
import PageNotFound from "@/pages/PageNotFound";
import { ReactQueryDevtools } from "@tanstack/react-query-devtools";
import { createBrowserRouter, RouterProvider } from "react-router-dom";
import AccountPage from "./pages/AccountPage";
import BudgetPage from "./pages/BudgetPage";
import TransactionPage from "./pages/TransactionPage";

function App() {
  const router = createBrowserRouter([
    {
      path: "/",
      element: <Home />,
    },
    {
      path: "/dashboard",
      element: <DashboardLayout />,
      children: [
        {
          index: true,
          element: <Dashboard />,
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
          loader: categoryLoader,
          errorElement: <PageError />,
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
    </>
  );
}

export default App;
