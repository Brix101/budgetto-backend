import * as z from "zod";

export const categorySchema = z.object({
  id: z.number(),
  created_at: z.string(),
  updated_at: z.string(),
  name: z.string(),
  note: z.string().nullish(),
});

export type Category = z.infer<typeof categorySchema>;

export const categoriesSchema = z.object({
  categories: z.array(categorySchema),
});
