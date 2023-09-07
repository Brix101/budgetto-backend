import { Row } from "@tanstack/react-table";

import { Icons } from "@/components/icons";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuShortcut,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { useBoundStore } from "@/lib/store";
import { Category } from "@/lib/validations/category";

interface DataTableRowActionsProps<TData> {
  row: Row<TData>;
}

export function CategoryDataTableRowActions<TData>({
  row,
}: DataTableRowActionsProps<TData>) {
  const { setMode } = useBoundStore((state) => state.category);
  const original = row.original as object;
  const isCreated = Object.keys(original).find((key) => key === "created_by");

  if (!isCreated) {
    return <></>;
  }

  function handleEditClick() {
    setMode({ mode: "update", category: original as Category });
  }

  function handleDeleteClick() {
    setMode({ mode: "delete", category: original as Category });
  }

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button
          variant="ghost"
          className="flex h-8 w-8 p-0 data-[state=open]:bg-muted"
        >
          <Icons.horizontalThreeDots className="h-4 w-4" />
          <span className="sr-only">Open menu</span>
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end" className="w-[160px]">
        <DropdownMenuItem onClick={handleEditClick}>Edit</DropdownMenuItem>
        <DropdownMenuSeparator />
        <DropdownMenuItem onClick={handleDeleteClick}>
          Delete
          <DropdownMenuShortcut></DropdownMenuShortcut>
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
