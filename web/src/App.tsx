import Home, { loader as homeLoader } from "@/pages/Home";
import PageError from "@/pages/PageError";
import SignIn from "@/pages/SignIn";
import { ReactQueryDevtools } from "@tanstack/react-query-devtools";
import { createBrowserRouter, RouterProvider } from "react-router-dom";
import CategoryPage from "./pages/CategoryPage";

function App() {
  const router = createBrowserRouter([
    {
      path: "/",
      element: <Home />,
      loader: homeLoader,
      errorElement: <PageError />,
      // children: [
      //   {
      //     path: "",
      //     element: <Team />,
      //     loader: teamLoader,
      //   },
      // ],
    },
    {
      path: "/sign-in",
      element: <SignIn />,
    },
    {
      path: "/categories",
      element: <CategoryPage />,
      loader: homeLoader,
    },
  ]);
  return (
    <>
      <RouterProvider router={router} />
      <ReactQueryDevtools initialIsOpen={false} position="bottom-right" />
    </>
  );
}

export default App;
