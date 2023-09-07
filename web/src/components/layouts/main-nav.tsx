import { Link } from "react-router-dom";

import { Icons } from "@/components/icons";

export function MainNav() {
  return (
    <div className="hidden gap-6 lg:flex">
      <Link
        aria-label="Home page"
        to="/"
        className="hidden items-center space-x-2 lg:flex"
      >
        <Icons.logo className="h-6 w-6" aria-hidden="true" />
        <span className="hidden font-bold lg:inline-block">Budgetto</span>
      </Link>
      {/* <NavigationMenu>
    <NavigationMenuList>
      {items?.[0]?.items ? (
        <NavigationMenuItem>
          <NavigationMenuTrigger className="h-auto">
            {items[0].title}
          </NavigationMenuTrigger>
          <NavigationMenuContent>
            <ul className="grid gap-3 p-6 md:w-[400px] lg:w-[500px] lg:grid-cols-[.75fr_1fr]">
              <li className="row-span-3">
                <NavigationMenuLink asChild>
                  <a
                    aria-label="Home"
                    className="flex h-full w-full select-none flex-col justify-end rounded-md bg-gradient-to-b from-muted/50 to-muted p-6 no-underline outline-none focus:shadow-md"
                    href="/"
                  >
                    <Icons.logo className="h-6 w-6" aria-hidden="true" />
                    <div className="mb-2 mt-4 text-lg font-medium">
                      Budgetto
                    </div>
                    <p className="text-sm leading-tight text-muted-foreground">
                      {siteConfig.description}
                    </p>
                  </a>
                </NavigationMenuLink>
              </li>
              {items[0].items.map((item) => (
                <ListItem
                  key={item.title}
                  title={item.title}
                  href={item.href}
                >
                  {item.description}
                </ListItem>
              ))}
            </ul>
          </NavigationMenuContent>
        </NavigationMenuItem>
      ) : null}
      
    </NavigationMenuList>
  </NavigationMenu> */}
    </div>
  );
}
