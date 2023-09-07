import { Auth0ContextInterface, User } from "@auth0/auth0-react";
import { Link } from "react-router-dom";

import { Icons } from "@/components/icons";
import { MainNav } from "@/components/layouts/main-nav";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Button, buttonVariants } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuShortcut,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Skeleton } from "@/components/ui/skeleton";
import { ModeToggle } from "../mode-toggle";

interface SiteHeaderProps {
  auth: Auth0ContextInterface<User>;
}

export function SiteHeader({ auth }: SiteHeaderProps) {
  const { user, logout, loginWithRedirect } = auth;

  const initials = `${user?.firstName?.charAt(0) ?? ""} ${
    user?.lastName?.charAt(0) ?? ""
  }`;

  return (
    <header className="sticky top-0 z-50 w-full border-b bg-background">
      <div className="container flex h-16 items-center">
        <MainNav />
        <div className="flex flex-1 items-center justify-end space-x-4">
          <nav className="flex items-center space-x-2">
            <ModeToggle />
            {auth.isLoading ? (
              <Skeleton className="h-10 w-14" />
            ) : (
              <>
                {user ? (
                  <DropdownMenu>
                    <DropdownMenuTrigger asChild>
                      <Button
                        variant="secondary"
                        className="relative h-8 w-8 rounded-full"
                      >
                        <Avatar className="h-8 w-8">
                          <AvatarImage
                            src={user.picture}
                            alt={user.nickname ?? ""}
                          />
                          <AvatarFallback>{initials}</AvatarFallback>
                        </Avatar>
                      </Button>
                    </DropdownMenuTrigger>
                    <DropdownMenuContent
                      className="w-56"
                      align="end"
                      forceMount
                    >
                      <DropdownMenuLabel className="font-normal">
                        <div className="flex flex-col space-y-1">
                          <p className="text-sm font-medium leading-none">
                            {user.name === user.email
                              ? user.nickname
                              : user.name}
                          </p>
                          <p className="text-xs leading-none text-muted-foreground">
                            {user.email}
                          </p>
                        </div>
                      </DropdownMenuLabel>
                      <DropdownMenuSeparator />
                      <DropdownMenuGroup>
                        <DropdownMenuItem asChild disabled className="hidden">
                          <Link to="/dashboard/profile">
                            <Icons.user
                              className="mr-2 h-4 w-4"
                              aria-hidden="true"
                            />
                            Profile
                            <DropdownMenuShortcut></DropdownMenuShortcut>
                          </Link>
                        </DropdownMenuItem>
                        <DropdownMenuItem asChild>
                          <Link to="/dashboard">
                            <Icons.terminal
                              className="mr-2 h-4 w-4"
                              aria-hidden="true"
                            />
                            Dashboard
                            <DropdownMenuShortcut></DropdownMenuShortcut>
                          </Link>
                        </DropdownMenuItem>
                        <DropdownMenuItem asChild disabled>
                          <Link to="/dashboard/settings">
                            <Icons.settings
                              className="mr-2 h-4 w-4"
                              aria-hidden="true"
                            />
                            Settings
                            <DropdownMenuShortcut></DropdownMenuShortcut>
                          </Link>
                        </DropdownMenuItem>
                      </DropdownMenuGroup>
                      <DropdownMenuSeparator />
                      <DropdownMenuItem asChild>
                        <a
                          className="w-full outline-none"
                          onClick={() =>
                            logout({
                              logoutParams: {
                                returnTo: window.location.origin,
                              },
                            })
                          }
                        >
                          <Icons.logout
                            className="mr-2 h-4 w-4"
                            aria-hidden="true"
                          />
                          Log out
                          <DropdownMenuShortcut></DropdownMenuShortcut>
                        </a>
                      </DropdownMenuItem>
                    </DropdownMenuContent>
                  </DropdownMenu>
                ) : (
                  <Button
                    className={buttonVariants({
                      size: "sm",
                    })}
                    onClick={() => loginWithRedirect()}
                  >
                    Sign In
                    <span className="sr-only">Sign In</span>
                  </Button>
                )}
              </>
            )}
          </nav>
        </div>
      </div>
    </header>
  );
}
