import { Button } from "@/components/ui/button";
import { useAuth0 } from "@auth0/auth0-react";
import { Link } from "react-router-dom";

function Home() {
  const { loginWithRedirect, isAuthenticated, logout, user, isLoading } =
    useAuth0();

  return (
    <div>
      {!isLoading ? (
        <div>
          {isAuthenticated && user ? (
            <div>
              <img src={user.picture} alt={user.name} />
              <h2>{user.name}</h2>
              <p>{user.email}</p>
            </div>
          ) : undefined}
          {isAuthenticated ? (
            <>
              <Button
                onClick={() =>
                  logout({ logoutParams: { returnTo: window.location.origin } })
                }
              >
                Logout
              </Button>
              <Link to={"/dashboard/categories"}>categories</Link>
            </>
          ) : (
            <Button onClick={() => loginWithRedirect()}>Login</Button>
          )}
        </div>
      ) : (
        <div>Loading...</div>
      )}
    </div>
  );
}

export default Home;
