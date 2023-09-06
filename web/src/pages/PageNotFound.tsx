import { Icons } from "@/components/icons";
import { Button } from "@/components/ui/button";
import { useNavigate } from "react-router-dom";

function PageNotFound() {
  const navigate = useNavigate();

  const homeClick = () => {
    navigate("/");
  };
  const backClick = () => {
    navigate(-1);
  };
  return (
    <section className="bg-white dark:bg-gray-900 ">
      <div className="container flex items-center min-h-screen px-6 py-12 mx-auto">
        <div>
          <p className="text-sm font-medium text-blue-500 dark:text-blue-400">
            404 error
          </p>
          <h1 className="mt-3 text-2xl font-semibold text-gray-800 dark:text-white md:text-3xl">
            We can’t find that page
          </h1>
          <p className="mt-4 text-gray-500 dark:text-gray-400">
            Sorry, the page you are looking for doesn't exist or has been moved.
          </p>

          <div className="flex items-center mt-6 gap-x-3">
            <Button variant="outline" onClick={backClick}>
              <Icons.moveLeft />
              <span>Go back</span>
            </Button>

            <Button onClick={homeClick}>Take me home</Button>
          </div>
        </div>
      </div>
    </section>
  );
}

export default PageNotFound;
