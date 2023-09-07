import { Category } from "@/lib/validations/category";
import { ColumnDef } from "@tanstack/react-table";

import { DataTableColumnHeader } from "@/components/data-table/data-table-column-header";
import { CategoryDataTableRowActions } from "@/components/data-table/data-table-row-actions.tsx";

export const categoryColumns: ColumnDef<Category>[] = [
  {
    accessorKey: "name",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Category" />
    ),
    cell: ({ row }) => (
      <div className="w-[100px] font-medium">{row.getValue("name")}</div>
    ),
  },
  {
    accessorKey: "note",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Description" />
    ),
    cell: ({ row }) => {
      return (
        <div className="flex space-x-2">
          <span className="max-w-[500px] truncate ">
            {row.getValue("note")}
          </span>
        </div>
      );
    },
    enableSorting: false,
    enableHiding: false,
  },
  {
    id: "actions",
    cell: ({ row }) => <CategoryDataTableRowActions row={row} />,
  },
];
