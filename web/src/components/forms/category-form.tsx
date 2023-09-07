import { QUERY_CATEGORIES_KEY } from "@/constant/query.constant";
import { Category, createCategorySchema } from "@/lib/validations/category";
import { createCategory } from "@/services/category.service";
import { useAuth0 } from "@auth0/auth0-react";
import { zodResolver } from "@hookform/resolvers/zod";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useForm } from "react-hook-form";
import * as z from "zod";

import { Icons } from "@/components/icons";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "@/components/ui/alert-dialog";
import { Button } from "@/components/ui/button";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Separator } from "@/components/ui/separator";
import { Textarea } from "@/components/ui/textarea";
import { useToast } from "@/components/ui/use-toast";
import { useBoundStore } from "@/lib/store";

type Inputs = z.infer<typeof createCategorySchema>;

export function CategoryCreateDialog() {
  const auth = useAuth0();
  const queryClient = useQueryClient();
  const { toast } = useToast();

  const { mode, setMode } = useBoundStore((state) => state.category);

  const form = useForm<Inputs>({
    resolver: zodResolver(createCategorySchema),
    defaultValues: {
      name: "",
      note: "",
    },
  });

  const { mutate, isLoading } = useMutation({
    mutationFn: createCategory,
    onSuccess: (category) => {
      queryClient.setQueriesData([QUERY_CATEGORIES_KEY], (prev: unknown) => {
        const categories = prev as Category[];
        return [category, ...categories];
      });

      setMode({ mode: "view" });
      toast({
        title: "Created successfully",
        description: `category ${category.name} created successfully`,
      });
    },
    onError: (error) => {
      console.log({ error });
    },
  });

  async function onSubmit(data: Inputs) {
    mutate({ auth, category: data });
  }

  function handleCancelClick() {
    setMode({ mode: "view" });
    form.reset();
  }

  return (
    <AlertDialog open={mode === "create"}>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Create category</AlertDialogTitle>
          <AlertDialogDescription>
            Create your new custom category
          </AlertDialogDescription>
        </AlertDialogHeader>
        <Separator />
        <div>
          <Form {...form}>
            <form
              className="grid gap-4"
              onSubmit={(...args) => void form.handleSubmit(onSubmit)(...args)}
            >
              <FormField
                control={form.control}
                name="name"
                disabled={isLoading}
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Name</FormLabel>
                    <FormControl>
                      <Input placeholder="category name" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="note"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Description</FormLabel>
                    <FormControl>
                      <Textarea
                        placeholder="category description"
                        className="resize-none"
                        {...field}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <AlertDialogFooter>
                <Button
                  type="button"
                  disabled={isLoading}
                  variant={"outline"}
                  onClick={handleCancelClick}
                >
                  Cancel
                </Button>
                <Button disabled={isLoading}>
                  {isLoading && (
                    <Icons.spinner
                      className="mr-2 h-4 w-4 animate-spin"
                      aria-hidden="true"
                    />
                  )}
                  Create
                  <span className="sr-only">Create</span>
                </Button>
              </AlertDialogFooter>
            </form>
          </Form>
        </div>
      </AlertDialogContent>
    </AlertDialog>
  );
}

export function CategoryDeleteDialog() {
  const { mode, category, setMode } = useBoundStore((state) => state.category);

  function handleCancelClick() {
    setMode({ mode: "view" });
  }

  if (category) {
    console.log({ mode, category });
  }
  return (
    <AlertDialog open={mode === "delete"}>
      <AlertDialogTrigger asChild>
        <Button variant="outline">Show Dialog</Button>
      </AlertDialogTrigger>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Are you absolutely sure?</AlertDialogTitle>
          <AlertDialogDescription>
            This action cannot be undone. This will permanently delete your
            cateory and remove your data from our servers.
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <Button variant={"destructive"}>Continue</Button>
          <AlertDialogCancel onClick={handleCancelClick}>
            Cancel
          </AlertDialogCancel>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  );
}

export function CategoryUpdateDialog() {
  const { mode, category, setMode } = useBoundStore((state) => state.category);

  function handleCancelClick() {
    setMode({ mode: "view" });
  }

  if (category) {
    console.log({ mode, category });
  }
  return (
    <AlertDialog open={mode === "update"}>
      <AlertDialogTrigger asChild>
        <Button variant="outline">Show Dialog</Button>
      </AlertDialogTrigger>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Are you absolutely sure?</AlertDialogTitle>
          <AlertDialogDescription>
            This action cannot be undone. This will permanently delete your
            cateory and remove your data from our servers.
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel onClick={handleCancelClick}>
            Cancel
          </AlertDialogCancel>
          <AlertDialogAction>Update</AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  );
}
