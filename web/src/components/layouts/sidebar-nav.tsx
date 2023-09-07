import type { SidebarNavItem } from "@/types";

import { Icons } from "@/components/icons";
import { cn } from "@/lib/utils";
import { Link, useLocation } from "react-router-dom";

export interface SidebarNavProps extends React.HTMLAttributes<HTMLDivElement> {
  items: SidebarNavItem[];
  isLoading?: boolean;
  loadingColor?: string;
}

export function SidebarNav({
  items,
  className,
  isLoading,
  loadingColor,
  ...props
}: SidebarNavProps) {
  const location = useLocation();

  if (!items?.length) return null;

  return (
    <div className={cn("flex w-full flex-col gap-2", className)} {...props}>
      {items.map((item, index) => {
        const Icon = Icons[item.icon ?? "chevronLeft"];

        if (isLoading) {
          return (
            <span
              className={cn(
                "h-9 w-full border border-transparent animate-pulse rounded-md bg-muted",
                loadingColor
              )}
            ></span>
          );
        }

        const isActive = item.href === location.pathname;
        return item.href ? (
          <Link
            aria-label={item.title}
            key={index}
            to={item.href}
            target={item.external ? "_blank" : ""}
            rel={item.external ? "noreferrer" : ""}
          >
            <span
              className={cn(
                "group flex w-full items-center rounded-md border border-transparent px-2 py-1 hover:bg-muted hover:text-foreground",
                isActive
                  ? "bg-muted font-medium text-foreground"
                  : "text-muted-foreground",
                item.disabled && "pointer-events-none opacity-60"
              )}
            >
              <Icon className="mr-2 h-4 w-4" aria-hidden="true" />
              <span>{item.title}</span>
            </span>
          </Link>
        ) : (
          <span
            key={index}
            className="flex w-full cursor-not-allowed items-center rounded-md p-2 text-muted-foreground hover:underline"
          >
            {item.title}
          </span>
        );
      })}
    </div>
  );
}
