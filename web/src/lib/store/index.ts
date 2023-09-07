import { create } from "zustand";
import { devtools } from "zustand/middleware";
import { CategorySlice, createBearSlice } from "./categorySlice";

type UnionToIntersection<U> = (
  U extends infer T ? (k: T) => void : never
) extends (k: infer I) => void
  ? I
  : never;

export type StoreState = UnionToIntersection<CategorySlice>;

const useBoundStore = create<StoreState>()(
  devtools((...a) => ({
    ...createBearSlice(...a),
  }))
);

export { useBoundStore };
