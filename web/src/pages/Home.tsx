import { SiteHeader } from "@/components/layouts/site-header";
import { useAuth0 } from "@auth0/auth0-react";

function Home() {
  const auth = useAuth0();

  return (
    <div>
      <SiteHeader auth={auth} />
    </div>
  );
}

export default Home;
