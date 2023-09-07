import { Category } from "@/lib/validations/category";
import { StateCreator } from "zustand";
import { StoreState } from ".";

type Mode = "view" | "create" | "update" | "delete";

type CategoryMode =
  | {
      mode: "update" | "delete";
      category: Category;
    }
  | {
      mode: "view" | "create";
      category?: Category;
    };

export interface CategorySlice {
  category: {
    mode: Mode;
    category?: Category;
    setMode: (mode: CategoryMode) => void;
  };
}

export const createBearSlice: StateCreator<
  StoreState,
  [],
  [],
  CategorySlice
> = (set) => ({
  category: {
    mode: "view",
    category: undefined,
    setMode: (mode) =>
      set(({ category }) => ({ category: { ...category, ...mode } })),
  },
});
