import { useRouteError } from "react-router-dom";

function PageError() {
  const error = useRouteError();
  console.log("++++++++++++++++++++++++++++++++++", { error });
  // Uncaught ReferenceError: path is not defined
  return <div>Dang!</div>;
}

export default PageError;
