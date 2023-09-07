import { Table } from "@tanstack/react-table";

import { Icons } from "@/components/icons";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";

interface DataTableToolbarProps<TData> {
  table: Table<TData>;
  placeHolder?: string;
}

export function DataTableToolbar<TData>({
  table,
  placeHolder,
}: DataTableToolbarProps<TData>) {
  const isFiltered = table.getState().columnFilters.length > 0;

  return (
    <div className="flex items-center justify-between">
      <div className="relative">
        <Icons.search className="absolute left-3 top-0 h-full w-5" />
        <Input
          placeholder={placeHolder ? placeHolder : "Filter..."}
          value={(table.getColumn("name")?.getFilterValue() as string) ?? ""}
          onChange={(event) =>
            table.getColumn("name")?.setFilterValue(event.target.value)
          }
          className="w-40 lg:w-80 px-10"
        />
        {isFiltered && (
          <Button
            variant="ghost"
            size={"icon"}
            onClick={() => table.resetColumnFilters()}
            className="absolute right-0 top-0 h-full px-3 py-1"
          >
            <Icons.close className="h-4 w-4" />
          </Button>
        )}
      </div>
    </div>
  );
}
